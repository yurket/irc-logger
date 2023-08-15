# IRC Logger

IRC Logger is a simple Golang application that logs IRC chat messages and provides a web interface to view the logs. It's dockerized and easily configurable.

## Usage

### Configuration

Set the following environment variables in the `docker-compose.yml`:

- `IRC_CHANNEL`: IRC channel to log
- `IRC_SERVER`: IRC server address
- `IRC_NICKNAME`: Bot's nickname
- `IRC_USERNAME`: IRC username
- `IRC_NAME`: Bot's name
- `IRC_TIMEZONE`: Timezone for logs (e.g., "Europe/London")
- `HTTP_USERNAME`: Username for HTTP authentication
- `HTTP_PASSWORD`: Password for HTTP authentication
- `HTTP_PORT`: HTTP port for the web interface

An example `docker-compose.yml.example` file is provided in the repository.


### Build and Run with Docker

```bash
docker-compose build
docker-compose up
```


### Access Logs

Access the logs at `http://localhost:8080`. You'll be prompted for the username and password specified in the environment variables.
