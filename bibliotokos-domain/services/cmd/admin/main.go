package main

import (
	"fmt"
	"os"
	"strconv"

	"bibliotokos.domain/pkg/datastore"
	"bibliotokos.domain/pkg/user"
)

func fatal(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)
	os.Exit(1)
}

func main() {
	args := os.Args[1:]
	if len(args) == 0 {
		fatal("usage: admin <db-path> user <subcommand> ...\n       admin database new <path>")
	}

	if args[0] == "database" {
		if len(args) < 3 || args[1] != "new" {
			fatal("usage: admin database new <path>")
		}
		ds, err := datastore.Open(args[2])
		if err != nil {
			fatal("error creating database: %v", err)
		}
		ds.Close()
		fmt.Println("database created:", args[2])
		return
	}

	if len(args) < 3 {
		fatal("usage: admin <db-path> user <subcommand> ...")
	}

	dbPath := args[0]
	if args[1] != "user" {
		fatal("unknown command %q", args[1])
	}

	ds, err := datastore.Open(dbPath)
	if err != nil {
		fatal("error opening database: %v", err)
	}
	defer ds.Close()

	uc, err := user.FromDB(user.Config{DB: ds})
	if err != nil {
		fatal("error initializing user controller: %v", err)
	}

	sub := args[2]
	rest := args[3:]

	switch sub {
	case "add":
		if len(rest) < 2 {
			fatal("usage: admin <db-path> user add <email> <password> [--is-admin]")
		}
		email := rest[0]
		password := rest[1]
		role := user.RoleStandard
		for _, a := range rest[2:] {
			if a == "--is-admin" {
				role = user.RoleAdmin
			}
		}
		id, err := uc.CreateUser(email, password, role)
		if err != nil {
			fatal("error creating user: %v", err)
		}
		fmt.Println(id)

	case "del":
		if len(rest) < 1 {
			fatal("usage: admin <db-path> user del <email|id>")
		}
		if err := uc.DeleteUser(rest[0]); err != nil {
			fatal("error deleting user: %v", err)
		}
		fmt.Println("deleted")

	case "password":
		if len(rest) < 2 {
			fatal("usage: admin <db-path> user password <email> <new-password>")
		}
		if err := uc.ResetPassword(rest[0], rest[1]); err != nil {
			fatal("error resetting password: %v", err)
		}
		fmt.Println("password updated")

	case "role":
		if len(rest) < 2 {
			fatal("usage: admin <db-path> user role [admin|standard] <email|id>")
		}
		var role user.Role
		switch rest[0] {
		case "admin":
			role = user.RoleAdmin
		case "standard":
			role = user.RoleStandard
		default:
			fatal("unknown role %q: must be admin or standard", rest[0])
		}
		if err := uc.SetRole(rest[1], role); err != nil {
			fatal("error setting role: %v", err)
		}
		fmt.Println("role updated")

	case "verify":
		if len(rest) < 1 {
			fatal("usage: admin <db-path> user verify <email|id>")
		}
		if err := uc.VerifyUser(rest[0]); err != nil {
			fatal("error verifying user: %v", err)
		}
		fmt.Println("user verified")

	case "list":
		offset, limit := 0, 50
		if len(rest) >= 1 {
			n, err := strconv.Atoi(rest[0])
			if err != nil {
				fatal("invalid offset %q", rest[0])
			}
			offset = n
		}
		if len(rest) >= 2 {
			n, err := strconv.Atoi(rest[1])
			if err != nil {
				fatal("invalid limit %q", rest[1])
			}
			limit = n
		}
		users, err := uc.ListUsers(offset, limit)
		if err != nil {
			fatal("error listing users: %v", err)
		}
		fmt.Printf("%-36s  %-30s  %-8s  %-8s  %s\n", "ID", "EMAIL", "ROLE", "VERIFIED", "CREATED")
		for _, u := range users {
			verified := "no"
			if u.Verified {
				verified = "yes"
			}
			fmt.Printf("%-36s  %-30s  %-8s  %-8s  %s\n",
				u.ID, u.Email, u.Role, verified, u.CreatedAt.Format("2006-01-02 15:04:05"),
			)
		}

	default:
		fatal("unknown user subcommand %q", sub)
	}
}
