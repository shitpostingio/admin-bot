db.createUser(
    {
        user: "automod",
        pwd: "automod",
        roles: [
            {
                role: "readWrite",
                db: "automod"
            }
        ]

    }
) 