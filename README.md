# startuplink-web-go ![Build Status](https://github.com/dlyahov/startuplink-web-go/workflows/build/badge.svg)

starutplink-web-go it's a lightweight web app written on Golang for saving your browser tabs in one place among browsers.

You can find this site [here](https://startuplink-web.com/home).

This app uses [Auth0](https://auth0.com/) for authentication users via OAuth2.

- [Install](#install)


## Install

This app can be installed using Docker. It's recommended way to install this app.

Original docker images can be found on [Dockerhub](https://hub.docker.com/repository/docker/aemdeveloper/startuplink-web-go).

### Parameters


| Environment variable    | Default                  | Description                                     |
| ----------------------- | ------------------------ | ----------------------------------------------- |
| AUTH0_CLIENT_ID         |                          | auth0 client id, required                       |
| AUTH0_CLIENT_SECRET     |                          | auth0 client secret, required                   |
| AUTH0_CALLBACK_URL      |                          | auth0 callback url, required                    |
| AUTH0_DOMAIN            |                          | auth0 personal domain, required                 |
| PORT                    | 8080                     | port of running application                     |
| STORE_BOLT_DB           | .                        | path where bolt db file willbe locatted         |
| TIMEOUT                 | 3000                     | timeout to connect to bolt db file              |
| COOKIE_AUTH_KEY         |                          | auth key to save cookie, required.
|                         |                          | It is recommended to use an authentication key with 32 or 64 bytes.|
| COOKIE_SECRET_KEY       |                          | The encryption key, if set, must be either 16, 24, or 32 bytes to select AES-128, AES-192, or AES-256 modes bytes.|
| PROFILE                 |                          | profile of application.Can have `local` or `prod` value for local and prod environment accordingly. required.                                      

### Register Auth0  

This app uses Auth0 because this resource allows you easy connect a lot of OAuth2 providers (and not only) relatively free.

You can find how generate your auth0 client for web app [here](https://auth0.com/docs/dashboard/guides/applications/register-app-regular-web). You can reuse it for your service.
