# senditsh

Senditsh is an secured oriented application allow user to send file from any where through a secured ssh connection. Receiver can verify origin of link and download file from the public link.
User Identification with jwt Token and github oauth token.

Stack:

    - Architecture: `monolithic`
    - Backend: `Go`.
    - Web framework: `Fiber`
    - Frontend: `Go` with `django` template.
    - Database: `Mongo`

## Run locally

Preparation: Create github oauth app and cliam github client id and github secret then fill in .env file. Example in example.env

1. Run mongodb service

   `make up`

   `docker-compose up`.

2. Config custom DNS for local developer.

   Example use `dnsmasq` with `dnsmasq.conf` config.

   - Copy dnsmasq.conf to /etc
   - start dnsmasq servie `systemctl start dnsmasq`

3. Config reversed proxy with nginx for sub domain feature.

   Example nginx config file `mysendit.conf`.

4. Start go service

   `make run`

TODO:

- [x] share file through secured connection.
- [x] authentication with github oauth service.
- [x] sub domain register for user.
- [ ] airdrop like feature. User share file to specific user, and wait for receiver accept for downloading.
