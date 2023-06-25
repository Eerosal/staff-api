# Staff API

Staff API is a simple JSON API that provides a list of players that have one of the configured groups on the Minecraft server.

Motivation behind this project was to provide a simple way to display a list of staff members on a website, but it can be used for other groups as well.

**NB:** This project has been done rapidly, so the software is provided "as is" and is not guaranteed to work in all environments. Use at your own risk. Feel free to open an issue if you find a bug.

## Data sources

The API uses the following data sources:
* LuckPerms MySQL (or MariaDB) database for fetching player groups
* Mojang API for fetching player names
* Avatars from any API that provides them in PNG format for a given UUID in the URL address (optional)

## API

### `GET /api/users`

Returns a list of users that have one of the configured groups on the server.

If avatars=true is passed in the query string, the response will contain a list of users with their avatars (base64 encoded).

#### Sample Response

```json
{
  "users": [
    {
      "uuid": "069a79f4-44e9-4726-a5be-fca90e38aaf5",
      "name": "Notch",
      "groups": [
        "admin"
      ],
      "avatar": "base64 encoded avatar if avatars=true is passed in the query string"
    },
    {
      "uuid": "61699b2e-d327-4a01-9f1e-0ea8c3f06bc6",
      "name": "Dinnerbone",
      "groups": [
        "moderator"
      ],
      "avatar": "base64 encoded avatar if avatars=true is passed in the query string"
    }
  ]
}
```

The API currently supports only one group per user. If a user has multiple groups, only the first one will be used.

## Setup

### Prerequisites

Docker and docker-compose are required to run this project. Preferably the latest versions.

### Configuration

By default, .env file is used for configuration. You can copy the example file (.env.example) and edit it to your liking.

List of available configuration options:
* `STAFF_API_REST_HTTP_ADDRESS` - Address on which the API will be listening for requests
* `STAFF_API_REST_HTTP_PORT` - Port on which the API will be listening for requests
* `STAFF_API_GROUPS` - Comma separated list of groups that will be used to fetch users. If * is passed, all groups will be used.
* `STAFF_API_IMAGE_URL` - URL address of the API that provides avatars for a given UUID. If empty, avatars will not be fetched. (%s will be replaced with the UUID)
* `STAFF_API_IMAGE_UPDATE_INTERVAL` - Interval in seconds between fetching new avatars from the API
* `STAFF_API_UPDATE_INTERVAL` - Interval in seconds between fetching new data from the data sources
* `STAFF_API_LUCKPERMS_CONNECTION_STRING` - Connection string to the LuckPerms MySQL database
* `STAFF_API_NAME_DATA_URL` - URL address of the Mojang API that provides player names for a given UUID. Mojang API is used by default. (%s will be replaced with the UUID)

It is highly recommended to use a reverse proxy (e.g. nginx) to handle SSL and rate limiting. 
The results are cacheable, so it is safe to use a high cache time.

The fetching of the name data is compliant with the rate limits of the Mojang API.

### Running

You can start the application just like any other docker-compose project or use the provided deploy.sh script.

### Testing

You can run the tests using the provided test.sh script. This will start the application with the required mock data sources and run the tests.

## Technologies used

* Golang
* Docker and docker-compose
* Wiremock
* MySQL
* Mojang API