package service

import (
	"context"
	"mime/multipart"
	"strconv"
	"strings"
	"sync"
	"time"

	google "github.com/google/uuid"

	"github.com/BangNopall/hology8-be/domain"
	"github.com/BangNopall/hology8-be/domain/contracts"
	"github.com/BangNopall/hology8-be/domain/dto"
	"github.com/BangNopall/hology8-be/domain/entity"
	"github.com/BangNopall/hology8-be/domain/enums"
	"github.com/BangNopall/hology8-be/pkg/aws"
	"github.com/BangNopall/hology8-be/pkg/helpers"
	"github.com/BangNopall/hology8-be/pkg/uuid"
)

type teamService struct {
	teamRepo        contracts.TeamRepository
	competitionRepo contracts.CompetitionRepository
	userRepo        contracts.UserRepository
	timeout         time.Duration
	aws             aws.CloudStorage
}

func NewTeamService(
	teamRepo contracts.TeamRepository,
	competitionRepo contracts.CompetitionRepository,
	userRepo contracts.UserRepository,
	timeout time.Duration,
	aws aws.CloudStorage,
) contracts.TeamService {
	return &teamService{
		teamRepo,
		competitionRepo,
		userRepo,
		timeout,
		aws,
	}
}

func (s *teamService) FetchTeamData(ctx context.Context, teamId string) (dto.TeamResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	res, _, err := s.teamRepo.FetchAllByConditionAndRelation(
		ctx,
		"id = ?",
		[]interface{}{teamId},
		"created_at ASC",
		nil,
		"Members",
		"Leader",
		"Members.User",
		"University",
		"Competition",
		"Announcements",
		"Competition.Announcements",
	)

	if err != nil {
		return dto.TeamResponse{}, err
	}

	if len(res) < 1 {
		return dto.TeamResponse{}, domain.ErrNotFound
	}

	team := dto.TeamEntityToResponse(&res[0])

	team.Leader = *dto.UserEntityToResponseDto(&res[0].Leader)

	team.Competition = dto.CompetitionEntityToResponse(&res[0].Competition)

	team.University = dto.UniversityEntityToDto(&res[0].University)

	team.Announcements = dto.AnnouncementSliceEntityToResponse(res[0].Announcements)

	team.Competition.Announcements = dto.AnnouncementSliceEntityToResponse(res[0].Competition.Announcements)

	select {
	case <-ctx.Done():
		return dto.TeamResponse{}, domain.ErrTimeout
	default:
		return team, nil
	}
}

func (s *teamService) FetchAll(
	ctx context.Context,
	params *dto.TeamParams,
	pageParam *dto.PaginationRequest,
) (dto.TeamPaginationResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var (
		condition   = ""
		args        = []interface{}{}
		order       = ""
		errChans    = make(chan error, 2)
		counterChan = make(chan dto.TeamCounter)
		teamsChan   = make(chan []entity.Team)
		pageChan    = make(chan dto.PaginationResponse)
		wg          sync.WaitGroup
	)

	if params.CompetitionID != 0 {
		condition = "competition_id = ?"
		args = []interface{}{params.CompetitionID}
	}

	if string(params.Status) != "" {
		condition = "status = ?"
		args = append(args, params.Status)
	}

	if string(params.Phase) != "" {
		condition = "phase = ?"
		args = append(args, params.Phase)
	}

	if string(params.WinnerPlace) != "" {
		condition = "winner_pace = ?"
		args = append(args, params.WinnerPlace)
	}

	if params.Name != "" {
		condition = "LOWER(team_name) LIKE LOWER(?)"
		args = append(args, "%"+params.Name+"%")
	}

	if params.SortBy == "latest" {
		order = "created_at DESC"
	} else {
		order = "created_at ASC"
	}

	wg.Add(2)

	go func() {
		defer wg.Done()

		counter, err := s.FetchStatisticData()

		if err != nil {
			errChans <- err
		}

		counterChan <- counter
	}()

	go func() {
		defer wg.Done()
		teams, pageResp, err := s.teamRepo.FetchAllByConditionAndRelation(
			ctx,
			condition,
			args,
			order,
			pageParam,
			"Members",
			"Leader",
			"Members.User",
			"University",
			"Competition",
			"Announcements",
		)

		if err != nil {
			errChans <- err
		}

		teamsChan <- teams
		pageChan <- pageResp
	}()

	go func() {
		wg.Wait()
		close(errChans)
		close(counterChan)
		close(teamsChan)
	}()

	teams := <-teamsChan
	pageResp := <-pageChan
	counter := <-counterChan

	for err := range errChans {
		if err != nil {
			if err == domain.ErrNotFound {
				return dto.TeamPaginationResponse{
					Teams:      []dto.TeamResponse{},
					Pagination: pageResp,
					Counter:    counter,
				}, nil
			}

			return dto.TeamPaginationResponse{
				Teams:      []dto.TeamResponse{},
				Pagination: pageResp,
				Counter:    counter,
			}, err
		}
	}

	res := make([]dto.TeamResponse, 0)

	for _, team := range teams {
		teamResp := dto.TeamEntityToResponse(&team)

		teamResp.Leader = *dto.UserEntityToResponseDto(&team.Leader)

		teamResp.Competition = dto.CompetitionEntityToResponse(&team.Competition)

		teamResp.University = dto.UniversityEntityToDto(&team.University)

		teamResp.Announcements = dto.AnnouncementSliceEntityToResponse(team.Announcements)

		res = append(res, teamResp)
	}

	resp := dto.TeamPaginationResponse{
		Teams:      res,
		Pagination: pageResp,
		Counter:    counter,
	}

	select {
	case <-ctx.Done():
		return dto.TeamPaginationResponse{}, domain.ErrTimeout
	default:
		return resp, nil
	}
}

func (s *teamService) FetchStatisticData() (dto.TeamCounter, error) {
	var (
		args      = []interface{}{"Verified", "Waiting", "Final", "Disqualified", "Eliminated"}
		condition = []string{"status = ?", "phase = ?"}
		res       = []int64{}
	)

	for idx, arg := range args {
		// cond idx detection
		if idx >= 2 {
			idx = 1
		} else {
			idx = 0
		}

		count, err := s.teamRepo.Count(condition[idx], []interface{}{arg}, "")

		if err != nil {
			return dto.TeamCounter{}, err
		}

		res = append(res, count)
	}

	counter := dto.TeamCounter{
		VerifiedTeam:     int(res[0]),
		OnholdTeam:       int(res[1]),
		FinalTeam:        int(res[2]),
		DisqualifiedTeam: int(res[3]),
		EliminatedTeam:   int(res[4]),
	}

	return counter, nil
}

func (s *teamService) FetchUserTeams(ctx context.Context, userId google.UUID) (dto.UserTeamsResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	team, _, err := s.teamRepo.FetchAllByConditionAndRelation(
		ctx,
		"leader_id = ?",
		[]interface{}{userId},
		"created_at ASC",
		nil,
		"Leader",
		"Competition",
		"Announcements",
	)

	if err != nil && err != domain.ErrNotFound {
		return dto.UserTeamsResponse{}, err
	}

	memberTeams, err := s.teamRepo.FetchMemberTeams(ctx, userId)

	if err != nil {
		return dto.UserTeamsResponse{}, err
	}

	res := dto.UserTeamsResponse{}

	if len(team) > 0 {
		teamResp := dto.TeamEntityToResponse(&team[0])

		teamResp.Leader = *dto.UserEntityToResponseDto(&team[0].Leader)

		teamResp.Competition = dto.CompetitionEntityToResponse(&team[0].Competition)

		teamResp.Announcements = dto.AnnouncementSliceEntityToResponse(team[0].Announcements)

		res.Teams = append(res.Teams, teamResp)
	}

	for _, memberTeam := range memberTeams {
		teamResp := dto.TeamEntityToResponse(&memberTeam.Team)

		teamResp.Leader = *dto.UserEntityToResponseDto(&memberTeam.Team.Leader)

		teamResp.Competition = dto.CompetitionEntityToResponse(&memberTeam.Team.Competition)

		teamResp.Announcements = dto.AnnouncementSliceEntityToResponse(memberTeam.Team.Announcements)

		res.Teams = append(res.Teams, teamResp)
	}

	select {
	case <-ctx.Done():
		return dto.UserTeamsResponse{}, domain.ErrTimeout
	default:
		return res, nil
	}
}

func (s *teamService) CreateTeam(ctx context.Context, leaderID google.UUID, team dto.TeamRegister) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	teamIDRandRune := helpers.GenerateRandomString(10)

	leaderIDPrefix := strings.Split(leaderID.String(), "-")[0]

	teamID := "H8" + strconv.Itoa(team.CompetitionID) + leaderIDPrefix + string(teamIDRandRune)

	competition, err := s.competitionRepo.FetchOneByID(ctx, team.CompetitionID)
	if err != nil {
		return err
	}

	competitionNameArray := strings.Split(competition.Name, " ")
	var competitionCode string
	for _, name := range competitionNameArray {
		competitionCode += string(name[0])
	}

	joinTokenRandRune := helpers.GenerateRandomString(15)

	team.JoinToken = competitionCode + "-" + joinTokenRandRune

	registeredTeams, err := s.teamRepo.FetchAllTeams(ctx, &dto.TeamParams{CompetitionID: team.CompetitionID}, "Members")
	if err == domain.ErrInternalServer {
		return err
	}

	if err == nil {
		for _, registeredTeams := range registeredTeams {
			for _, member := range registeredTeams.Members {
				if member.UserID == leaderID {
					return domain.ErrUserAlreadyRegistered
				}
			}
		}
	}

	var user entity.User

	err = s.userRepo.FindUser(&user, &dto.UserParam{ID: leaderID})

	if err != nil {
		return err
	}

	newTeam := entity.NewTeam(teamID, leaderID, team.CompetitionID, team.Name, team.JoinToken, enums.Waiting, enums.Elimination, enums.Default)

	newTeam.UniversityID = user.UniversityID

	err = s.teamRepo.InsertTeam(ctx, newTeam)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *teamService) UploadPaymentProof(ctx context.Context, id string, userId string, paymentFile *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if paymentFile.Size > 1024*1024 {
		return domain.ErrFileTooBig
	}

	team, err := s.teamRepo.FetchOneByParams(ctx, &dto.TeamParams{ID: id})

	if err != nil {
		return err
	}

	if userId != team.LeaderID.String() {
		return domain.ErrUserNotTeamLeader
	}

	var link string

	path := "teams/" + id + "/payment-proof/" + uuid.New().String() + " " + time.Now().Format("02-Jan-2006 15:04:05")

	if team.PaymentProofLink != "" {

		link, err = s.aws.Update(path, paymentFile, team.PaymentProofLink)

	} else {

		link, err = s.aws.Upload(path, paymentFile)

	}

	if err != nil {
		return err
	}

	err = s.teamRepo.UpdateTeam(ctx, id, &dto.TeamUpdate{PaymentProofLink: link})

	if err != nil {
		return err
	}

	return nil
}

func (s *teamService) UploadTwibbonProof(ctx context.Context, id string, userId string, twibbonFile *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if twibbonFile.Size > 1024*1024 {
		return domain.ErrFileTooBig
	}

	team, err := s.teamRepo.FetchOneByParams(ctx, &dto.TeamParams{ID: id})

	if err != nil {
		return err
	}

	if userId != team.LeaderID.String() {
		return domain.ErrUserNotTeamLeader
	}

	var link string

	path := "teams/" + id + "/twibbon-proof/" + uuid.New().String() + " " + time.Now().Format("02-Jan-2006 15:04:05")

	if team.TwibbonProofLink != "" {

		link, err = s.aws.Update(path, twibbonFile, team.TwibbonProofLink)

	} else {

		link, err = s.aws.Upload(path, twibbonFile)

	}

	if err != nil {
		return err
	}

	err = s.teamRepo.UpdateTeam(ctx, id, &dto.TeamUpdate{TwibbonProofLink: link})

	if err != nil {
		return err
	}

	return nil
}

func (s *teamService) UploadProposalDoc(ctx context.Context, id string, userId string, proposalFile *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if proposalFile.Size > 100*1024*1024 {
		return domain.ErrFileTooBig
	}

	team, err := s.teamRepo.FetchOneByParams(ctx, &dto.TeamParams{ID: id})

	if err != nil {
		return err
	}

	if userId != team.LeaderID.String() {
		return domain.ErrUserNotTeamLeader
	}

	var link string

	path := "teams/" + id + "/proposal-doc/" + uuid.New().String() + " " + time.Now().Format("02-Jan-2006 15:04:05")

	if team.ProposalDocLink != "" {

		link, err = s.aws.Update(path, proposalFile, team.ProposalDocLink)

	} else {

		link, err = s.aws.Upload(path, proposalFile)

	}

	if err != nil {
		return err
	}

	err = s.teamRepo.UpdateTeam(ctx, id, &dto.TeamUpdate{ProposalDocLink: link})

	if err != nil {
		return err
	}

	return nil
}

func (s *teamService) UploadStatementLetter(ctx context.Context, id string, userId string, statementLetter *multipart.FileHeader) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if statementLetter.Size > 100*1024*1024 {
		return domain.ErrFileTooBig
	}

	team, err := s.teamRepo.FetchOneByParams(ctx, &dto.TeamParams{ID: id})

	if err != nil {
		return err
	}

	if userId != team.LeaderID.String() {
		return domain.ErrUserNotTeamLeader
	}

	var link string

	path := "teams/" + id + "/statement-letter/" + uuid.New().String() + " " + time.Now().Format("02-Jan-2006 15:04:05")

	if team.StatementLetterLink != "" {

		link, err = s.aws.Update(path, statementLetter, team.StatementLetterLink)

	} else {

		link, err = s.aws.Upload(path, statementLetter)

	}

	if err != nil {
		return err
	}

	err = s.teamRepo.UpdateTeam(ctx, id, &dto.TeamUpdate{StatementLetterLink: link})

	if err != nil {
		return err
	}

	return nil
}

func (s *teamService) UpdateTeamData(ctx context.Context, id string, userId string, team *dto.TeamUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	if string(team.Status) != "" || string(team.Phase) != "" || string(team.WinnerPlace) != "" {
		return domain.ErrForbiddenUpdate
	}

	if team.LeaderID != google.Nil {
		return domain.ErrForbiddenUpdate
	}

	teamEntity, err := s.teamRepo.FetchOneByID(ctx, id)
	if err != nil {
		return domain.ErrNotFound
	}

	if userId != teamEntity.LeaderID.String() {
		return domain.ErrUserNotTeamLeader
	}

	err = s.teamRepo.UpdateTeam(ctx, id, team)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *teamService) UpdateTeamStatus(ctx context.Context, id string, team *dto.TeamUpdate) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.teamRepo.UpdateTeam(ctx, id, team)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *teamService) UpdateLeader(ctx context.Context, teamId string, leaderId google.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	leaderTeam, err := s.teamRepo.FetchOneByParams(ctx, &dto.TeamParams{ID: teamId})

	if err != nil {
		return err
	}

	teams, err := s.teamRepo.FetchAllTeams(ctx, &dto.TeamParams{CompetitionID: leaderTeam.CompetitionID}, "Members")

	if err != nil {
		return err
	}

	// check whether user already registered to a team with same competition
	for _, team := range teams {
		for _, member := range team.Members {
			if member.UserID == leaderId && member.TeamID != teamId {
				return domain.ErrUserAlreadyRegistered
			}
		}
	}

	var (
		errChans = make(chan error, 1)
		wg       sync.WaitGroup
	)

	wg.Add(2)

	s.teamRepo.Begin()

	go func() {
		defer wg.Done()
		err := s.teamRepo.UpdateTeam(ctx, teamId, &dto.TeamUpdate{LeaderID: leaderId})

		if err != nil {
			errChans <- err
		}
	}()

	go func() {
		defer wg.Done()
		s.teamRepo.DeleteTeamMember(ctx, &entity.DetailTeams{UserID: leaderId, TeamID: teamId})
	}()

	go func() {
		wg.Wait()
		close(errChans)
	}()

	for err := range errChans {
		if err != nil {
			s.teamRepo.Rollback()
			return err
		}
	}

	select {
	case <-ctx.Done():
		s.teamRepo.Rollback()
		return domain.ErrTimeout
	default:
		s.teamRepo.Commit()
		return err
	}
}

func (s *teamService) RemoveMember(ctx context.Context, member *entity.DetailTeams) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	err := s.teamRepo.DeleteTeamMember(ctx, member)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *teamService) JoinTeam(ctx context.Context, joinToken string, userID google.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	var user entity.User
	err := s.userRepo.FindUser(&user, &dto.UserParam{ID: userID}, "Teams", "Teams.Team")

	if err != nil {
		return err
	}

	team, err := s.teamRepo.FetchOneByParams(ctx, &dto.TeamParams{JoinToken: joinToken}, "Members")

	if err != nil {
		return err
	}

	uTL, err := s.teamRepo.FetchOneByParams(ctx, &dto.TeamParams{LeaderID: userID})

	if err == nil && uTL.CompetitionID == team.CompetitionID {
		return domain.ErrUserAlreadyRegistered
	}

	for _, t := range user.Teams {
		// check if user already registered to a team with the same competition
		if t.Team.CompetitionID == team.CompetitionID {
			return domain.ErrUserAlreadyRegistered
		}
	}

	// check if a team is already contains 2 members (exc leader)
	if len(team.Members) >= 2 {
		return domain.ErrTeamFull
	}

	detailTeam := entity.DetailTeams{
		UserID: userID,
		TeamID: team.ID,
	}

	err = s.teamRepo.InsertTeamMember(ctx, &detailTeam)

	select {
	case <-ctx.Done():
		return domain.ErrTimeout
	default:
		return err
	}
}

func (s *teamService) CountTeamNUniv(ctx context.Context) (dto.TeamNUnivCounter, error) {
	ctx, cancel := context.WithTimeout(ctx, s.timeout)
	defer cancel()

	totalTeams, err := s.teamRepo.Count("1 = 1", nil, "")
	if err != nil {
		return dto.TeamNUnivCounter{}, err
	}

	totalUniv, err := s.teamRepo.Count("1 = 1", nil, "university_id")

	teamNUnivCounter := dto.TeamNUnivCounter{
		TeamCounter: int(totalTeams),
		UnivCounter: int(totalUniv),
	}

	select {
	case <-ctx.Done():
		return teamNUnivCounter, domain.ErrTimeout
	default:
		return teamNUnivCounter, err
	}
}
