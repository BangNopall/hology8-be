async function fetchData() {
    const res = await fetch("http://localhost:8080/api/v1/users/oauth", {
        method: 'GET',
        headers: {
            "Content-Type": "application/json",
            "x-api-key": "Key 987adb66d54716e09086d045ca683f4aea45702067785df61c631ade1d62d9f7"
        }
    })

    const data = await res.json()

    if (data.data.redirect_link) {
        window.location.href = data.data.redirect_link
    }
}

fetchData()