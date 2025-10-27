# FoodFinder

A tool for University of Minnesota students to get notifications for their
favorite dining hall meals. This project is associated with
[Social Coding](https://www.socialcoding.net/). If you are in the group and need
access to the repository, email your github username to `calle159@umn.edu`.

## Frontend

### Stack

-   Package Manager - [Bun](https://bun.com/)
-   Framework - [Next.js](https://nextjs.org/)
-   Styling - [Tailwind CSS](https://tailwindcss.com/)
-   Formatter - [Prettier](https://prettier.io/)
-   Linter - [Eslint](https://eslint.org/)

### Getting started

-   Install bun.
-   Run `bun install` to install all required packages.
-   Run `bun pm trust --all` to trust installed packages.
-   Run `bun dev` to start the dev server.

### Before you push

-   Run `bun lint` to lint your code for common issues.
-   Run `bun format` to format your code.

## Backend

### Stack

-   Language [Go](https://go.dev/)
-   API [Gin](https://github.com/gin-gonic/gin)

### Environment Variables

#### JWT
Create a `.env` file in the root of the repository with the following values:

-   access_key=YOUR-KEY
-   refresh_key=YOUR-KEY

#### jsonv2

The `dineocclient` package uses the `encodings/json/v2` API. Currently, this
library is not available unless you set the following environment variable prior
to builds: `GOEXPERIMENT=jsonv2`

To set this environment variable on powershell, use the following command:

`$env:GOEXPERIMENT="jsonv2"`

Using this command will set the environment variable for the rest of your
powreshell session.

In bash, you can use the following command:

`export GOEXPERIMENT=jsonv2`

Other *nix shells will have similar syntax, if you aren't using bash then you can
look up specific methods to set variables in your particular shell.

### Notifier
The module `notifier` located at `/notifier` from the root of the project
reads its sending address as `NOTIFIER_EMAIL`, its SMTP password as
`NOTIFIER_PASSWORD`, and the database connection string as `DATABASE_URL`. The
notifier should not generally be run from dev machines without prior discussion.
However, you must make sure to set these environment variables before running
it or else you will get an error.

Make sure run `go get .` in the `notifier` directory as well.

### Getting started

-   Install Go.
-   Run `go get .` to install packages.
-   Run `go install github.com/air-verse/air@latest` to install air.
-   Run `air` to start the dev server.

### Before you push

-   Run `go fmt .` format your code.
