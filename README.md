# Anbox cloud management


## Prerequisite : Load seed database

* First, download the tools :

```
cd scripts/postgresql && \
    make psql-add-postgres-ubuntu-client && \
    make psql-add-migrate-ubuntu-client
```

* Then, launch the database :

```
cd scripts/postgresql && \
    make psql-run
```

* Apply the migrations :

```
cd scripts/postgresql && \
    make psql-migrate-up
```

(Check that you have a PostgreSQL Docker container running : `docker ps | grep "anboxcloud_postgresql_server"`)

* (Optional) Seed the database :

```
cd scripts/postgresql && \
    make psql-import-db
```

## Building the server and the CLI

* To build the binaries, simply do :

```
make build
```

* To run the server, do :

```
make run
```

* The CLI client to interact with the server is at `bin/anbox-cli` (**use the full `./bin/anbox-client` path when executing, else some ENV variables won't be declared and the client will panic**):

```
$ ./bin/anbox-cli --help

Managing the Anbox application. We can do CRUD operations on "users" and "games" and also create links between entities

Usage:
  anbox-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create entities
  delete      Delete entities
  help        Help about any command
  link        Link entities
  list        List entities
  login       Login to a user account
  update      Update entities

Flags:
  -h, --help   help for anbox-cli

Use "anbox-cli [command] --help" for more information about a command.
```

## OpenAPI/Postman integration

* Once you have a running server, you can import the **OpenAPI** specification at `pkg/api/open-api.yml` in Postman to directly interact with the server through the beautiful Postman UI. This solution is ideal to see the API  documentation and use the prefilled Postman queries for each API calls.

## Fun features

* `user` and `game` CRUD
* `metadata` association
* Random game traffic simulator in `pkg/server/server.go`
* JWT auth
* OpenAPI integration
* Migration tooling is quite efficient
* Highly flexible and helping CLI 

## Pro tips

* If you use the CLI, you can first create a user (no auth allowed) and then login with it to be allowed to interact with the system :

```
./bin/anbox-cli create user --email=gabriel.mougard123@gmail.com \
  --username=gabrielmougard123 \
  --password=gabrielmougard123

...


./bin/anbox-cli login --email=gabriel.mougard123@gmail.com --password=gabrielmougard123
```

After that you have a JWT on your disk that will allow all the calls ! The CLI mirror the OpenAPI spec (and the `cobra` CLI is also quite helping)

* You can add games (`./bin/anbox-cli create game [FLAGS]`), link a game to a user (`./bin/anbox-cli link game [FLAGS]`), list the users (with their metadata !), example:

```
$ ./bin/anbox-cli list user
{
	"usersCount": 14,
	"usersWithMetadata": [
		{
			"user": {
				"email": "gabriel.mougard@gmail.com",
				"age": 24,
				"username": "root",
				"createdAt": "2022-10-27T03:55:03.281924Z",
				"updatedAt": "2022-10-27T03:55:03.281924Z"
			},
			"metadata": [
				{
					"player": {
						"email": "gabriel.mougard@gmail.com",
						"age": 24,
						"username": "root",
						"createdAt": "2022-10-27T03:55:03.281924Z",
						"updatedAt": "2022-10-27T03:55:03.281924Z"
					},
					"playedGame": {
						"title": "leagueoflegends",
						"description": "League of Legends: Wild Rift is a 5v5 MOBA game.",
						"url": "https://wildrift.leagueoflegends.com/en-us/",
						"ageRating": 16,
						"publisher": "riotgames",
						"createdAt": "2022-10-27T06:24:47.950425Z",
						"updatedAt": "2022-10-27T06:24:47.950425Z"
					},
					"playTime": 7996,
					"playTimeHuman": "133 h 16 min",
					"playerUsername": "root",
					"gameTitle": "leagueoflegends",
					"createdAt": "2022-10-27T07:27:40.442035Z",
					"updatedAt": "2022-10-27T09:13:41.823645Z"
				},
				{
					"player": {
						"email": "gabriel.mougard@gmail.com",
						"age": 24,
						"username": "root",
						"createdAt": "2022-10-27T03:55:03.281924Z",
						"updatedAt": "2022-10-27T03:55:03.281924Z"
					},
					"playedGame": {
						"title": "pokemon go",
						"description": "Play with pokemon on reaL LIFE WITH ANDROID !",
						"url": "https://pokemongolive.com/",
						"ageRating": 7,
						"publisher": "nintendo",
						"createdAt": "2022-10-27T06:31:06.357362Z",
						"updatedAt": "2022-10-27T06:31:06.357362Z"
					},
					"playTime": 6658,
					"playTimeHuman": "110 h 58 min",
					"playerUsername": "root",
					"gameTitle": "pokemon go",
					"createdAt": "2022-10-27T07:26:58.33622Z",
					"updatedAt": "2022-10-27T09:13:41.824862Z"
				}
			]
		},
		{
			"user": {
				"email": "gabriel.mougard2@gmail.com",
				"age": 1,
				"username": "root2",
				"createdAt": "2022-10-27T03:56:36.61615Z",
				"updatedAt": "2022-10-27T03:56:36.61615Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "jean.paul@gmail.com",
				"age": 24,
				"username": "jeanpaul",
				"createdAt": "2022-10-27T04:27:38.915887Z",
				"updatedAt": "2022-10-27T04:27:38.915887Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "meera.rowley@gmail.com",
				"age": 12,
				"username": "meerarowley",
				"createdAt": "2022-10-27T06:13:41.984611Z",
				"updatedAt": "2022-10-27T06:13:41.984611Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "jarvis.reid@gmail.com",
				"age": 16,
				"username": "jarvisreid",
				"createdAt": "2022-10-27T06:14:20.767707Z",
				"updatedAt": "2022-10-27T06:14:20.767707Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "bradenfeeney@gmail.com",
				"age": 6,
				"username": "bradenfeeney",
				"createdAt": "2022-10-27T06:14:50.29954Z",
				"updatedAt": "2022-10-27T06:14:50.29954Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "clarencevazquez@gmail.com",
				"age": 33,
				"username": "clarencevazquez",
				"createdAt": "2022-10-27T06:15:32.673177Z",
				"updatedAt": "2022-10-27T06:15:32.673177Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "ritawest@gmail.com",
				"age": 14,
				"username": "ritawest",
				"createdAt": "2022-10-27T06:16:07.053839Z",
				"updatedAt": "2022-10-27T06:16:07.053839Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "alyssiamonaghan@gmail.com",
				"age": 15,
				"username": "alyssiamonaghan",
				"createdAt": "2022-10-27T06:16:31.464686Z",
				"updatedAt": "2022-10-27T06:16:31.464686Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "jovanpatrick@gmail.com",
				"age": 23,
				"username": "jovanpatrick",
				"createdAt": "2022-10-27T06:16:57.995594Z",
				"updatedAt": "2022-10-27T06:37:26.424854Z"
			},
			"metadata": [
				{
					"player": {
						"email": "jovanpatrick@gmail.com",
						"age": 23,
						"username": "jovanpatrick",
						"createdAt": "2022-10-27T06:16:57.995594Z",
						"updatedAt": "2022-10-27T06:37:26.424854Z"
					},
					"playedGame": {
						"title": "callofduty",
						"description": "modified callofduty description",
						"url": "https://modifiedcallofdutyurl.com",
						"ageRating": 21,
						"publisher": "activision",
						"createdAt": "2022-10-27T06:25:45.697334Z",
						"updatedAt": "2022-10-27T06:28:01.417843Z"
					},
					"playTime": 9219,
					"playTimeHuman": "153 h 39 min",
					"playerUsername": "jovanpatrick",
					"gameTitle": "callofduty",
					"createdAt": "2022-10-27T06:51:07.33468Z",
					"updatedAt": "2022-10-27T09:13:41.825713Z"
				}
			]
		},
		{
			"user": {
				"email": "rhiannelson@gmail.com",
				"age": 20,
				"username": "rhiannelson",
				"createdAt": "2022-10-27T06:17:21.574769Z",
				"updatedAt": "2022-10-27T06:17:21.574769Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "eliotttalbot@gmail.com",
				"age": 21,
				"username": "eliotttalbot",
				"createdAt": "2022-10-27T06:17:44.41983Z",
				"updatedAt": "2022-10-27T06:17:44.41983Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "lyndonstubbs@gmail.com",
				"age": 7,
				"username": "lyndonstubbs",
				"createdAt": "2022-10-27T06:18:07.442724Z",
				"updatedAt": "2022-10-27T06:18:07.442724Z"
			},
			"metadata": []
		},
		{
			"user": {
				"email": "truc.truc@gmail.com",
				"age": 2,
				"username": "tructruc123",
				"createdAt": "2022-10-27T09:09:12.526627Z",
				"updatedAt": "2022-10-27T09:09:12.526627Z"
			},
			"metadata": []
		}
	]
}
```

## Conclusion

I took a bit more time (around 12h) than expected but the project looks really nice now !
