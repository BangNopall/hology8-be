package flag

import (
	"flag"
)

type Flag struct {
	Fresh       bool
	Seeder      bool
	SeederModel string
}

var FlagVars = getFlags()

func getFlags() *Flag {
	fresh := flag.Bool("fresh", false, "Dropping all database tables before running new migration")
	seeder := flag.Bool("seeder", false, "Inserting all seeders to database")
	seederModel := flag.String("model", "", "Specify seeder model to run\nMust be used with flag seeder\nAvailable models can be looked up at pgsql_conn.go file func Seeder")

	flag.Parse()

	flag := &Flag{
		Fresh:       *fresh,
		Seeder:      *seeder,
		SeederModel: *seederModel,
	}

	return flag
}
