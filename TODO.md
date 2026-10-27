# TODO — KaellyRetro launch

Tick items as they are done. Scope: `/about`, `/help`, `/align`, `/job` and `/config` (`get`, `rss`, `server`, `twitter`), all tagged `DOFUS_RETRO` (`game = 3`).

## Code (this repo)

- [x] Strip kaelly-discord down to the Retro scope
- [x] `constants.Name = "KaellyRetro"`, so the reply queue `KaellyRetro-<shard>.answers` doesn't collide with KaellyBot's
- [x] Job level capped at 100, alignment level capped at 60
- [x] Bot avatar and neutral city image loaded from `kaelly-cdn/kaellyretro/`
- [x] `/config rss` and `/config twitter` restored from kaelly-discord (no almanax)
- [x] Tutorial links point to `kaelly-cdn/kaellyretro/tutorials/`
- [x] CI: `RETRO_` secrets for what differs from KaellyBot (token, shard count, command IDs), KaellyBot's own secrets for the rest
- [x] kaelly-commands upgraded to v1.1.1
- [x] First commit
- [ ] Decide which bot to recommend for unsupported features, and where to mention it (welcome, `/about`, `/help`), in the 5 languages
- [ ] Upload the tutorial screenshots to `kaelly-cdn/kaellyretro/tutorials/`: `help.webp`, then `<name>_<LANG>.webp` (LANG = DE, EN, ES, FR, PT) for `about`, `align_get`, `align_set`, `config_get`, `config_rss`, `config_server`, `config_twitter`, `job_get`, `job_set`
- [x] Help examples: `Devblog` for `/config rss`, the Retro accounts for `/config twitter` (`DofusRETRO_EN` for EN and DE, `DOFUSRetro_FR`, `DOFUSRetro_ES`, `DOFUSRetro_PTBR`)
- [x] Boune as the example server in the help texts (a Retro server on the CDN)

## Discord applications

- [x] Create the Retro Discord application (production)
- [ ] Development application, if needed
- [ ] Upload the emojis to each Retro application:
  - [x] Cities
  - [x] Orders
  - [x] Jobs
  - [x] Servers
  - [ ] Miscellaneous

## kaelly-commands / kaelly-registrar

- [x] kaelly-commands v1.1.0: `GetRetroCommands()`
- [x] kaelly-commands v1.1.1: Retro alignment capped at 60, `/config rss` and `/config twitter` for Retro
- [x] `/config` description for Retro: the registrar's `config.description` no longer mentions Almanax, and fits Discord's 100 characters
- [x] Registrar: `GAME` (1/2/3) picks the command set, kaelly-commands v1.1.1, second Helm release `kaelly-registrar-retro`
- [x] Registrar: no `RETRO_CLIENT_ID` secret, it runs locally with `GAME=3`, `CLIENT_ID=1239972754132832256`, `TOKEN=<retro token>`, `PRODUCTION=true`
- [x] Register the Retro commands and copy the returned IDs into the `RETRO_ABOUT_ID`, `RETRO_ALIGN_ID`, `RETRO_CONFIG_ID`, `RETRO_HELP_ID` and `RETRO_JOB_ID` secrets

## Data

Labels use `locale` FR=1, EN=2, ES=3, DE=4, PT=5. EN is the fallback, so it is mandatory.

Icons and images are on the CDN under `https://raw.githubusercontent.com/kaellybot/kaelly-cdn/refs/heads/main/kaellyretro/`: `cities/` (bonta, brakmar), `jobs/` (33), `servers/icons/` and `servers/images/` (allisteria, boune, fallanster).

- [x] `servers` + `server_labels` (`game = 3`, `icon`, `image`)
- [x] `jobs` + `job_labels` (`game = 3`, `icon`, `color`)
- [x] `cities` + `city_labels` (`game = 3`, `icon`, `type` light/dark, `color`)
- [x] `orders` + `order_labels` (`game = 3`)
- [x] `feed_sources` (`game = 3`, `feed_type_id`, `locale`, `url`, `news_channel_id`): feed types and their labels are shared with DOFUS
- [x] `twitter_accounts` (`game = 3`, `id`, `name`, `locale`, `news_channel_id`)
- [ ] `emojis` rows for Retro IDs unknown to Dofus (the table is shared, it has no `game` column)
- [ ] `emoji_snowflakes` (`game = 3`, `production` 1 and/or 0) for:
  - [ ] `Server`: each server
  - [ ] `Job`: each job
  - [ ] `City`: each city
  - [ ] `Order`: `light<order>`, `dark<order>` and `<order>` for each order
  - [ ] `Miscellaneous`: `first`, `previous`, `next`, `last`, `global`, `scroll`, `github`, `twitter`, `rss`, `dofusGame` (Retro logo)

## Other services

- [x] kaelly-books: accept `DOFUS_RETRO` (BOOK-2 only accepted `DOFUS_GAME` and `DOFUS_TOUCH`), otherwise `/align` and `/job` fail (deployed)
- [x] kaelly-rss / kaelly-twitter: checked, Retro rows are read and published as `DOFUS_RETRO` like the other games
- [x] kaelly-notifier: serve Retro (`DISCORD_TOKEN_DOFUS_RETRO`, bot `KaellyRetro`)
- [x] kaelly-notifier: its CI reads `RETRO_DISCORD_TOKEN`, which must exist **before** deploying: the notifier refuses to start without a served game's token, Dofus included. Retro news published until then is lost
- [x] kaelly-notifier: restart after inserting Retro data, it loads `feed_sources` and `twitter_accounts` only at startup
- [ ] Retro application: member of the server hosting the Retro news channels (post and crosspost), and of the one hosting the shared reporting channel (`REPORTING_CHANNEL_ID`), where it posts its join/leave reports
- [x] kaelly-configurator: accept `DOFUS_RETRO` (`isGameSupported` only accepted `DOFUS_GAME` and `DOFUS_TOUCH`), otherwise every `/config` write and guild join fails (deployed)
- [ ] Dashboards: group or filter by the `game` tag (`DOFUS_RETRO`), so Retro requests aren't mixed with Dofus ones

## Deployment

- [x] Create the Retro secrets: `RETRO_DISCORD_TOKEN`, `RETRO_DISCORD_SHARD_COUNT` (plus the command IDs above). MySQL, RabbitMQ, ports, log level and production reuse KaellyBot's secrets
- [ ] First release (the Helm deployment runs on tags)
- [ ] After deploy: no "No entity emoji found" / "No miscellaneous emoji found" warnings in the logs

## Documentation

- [ ] kaelly-amqp `docs/specs/game-isolation.md` §4.10: kaelly-commands is shared (one function per game) instead of forked
