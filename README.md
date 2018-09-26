
# Running

go install ./...; source /Users/perelin/.zshrc

OR

make

# Logs

https://apps.eu.sematext.com/ 

# Notes

https://mholt.github.io/json-to-go/

https://stackoverflow.com/questions/28582728/running-heroku-rundetached-programmatically-how-exactly
https://devcenter.heroku.com/articles/one-off-dynos#running-tasks-in-background

heroku run:detached recruitbot3000
heroku logs --app secret-waters-77818 --dyno run.4980 --tail

heroku run:detached jb3k-cr-ads-monster
heroku run jb3k-tools-db

https://ieftimov.com/golang-package-multiple-binaries

go install ./... 

https://stackoverflow.com/questions/24790175/when-is-the-init-function-run

https://cloud.google.com/natural-language/
https://azure.microsoft.com/en-us/services/cognitive-services/text-analytics/ 
https://www.textrazor.com/plans

https://benjamincongdon.me/blog/2018/03/01/Scraping-the-Web-in-Golang-with-Colly-and-Goquery/

https://stackoverflow.com/questions/42386975/heroku-postgresql-with-google-datastudio

INteresting: https://www.scraperapi.com/pricing

http://www.postgresqltutorial.com/postgresql-rename-column/


# hot reload

$ gin

OR

$ fresh

# heroku pg backups
$ heroku pg:backups:capture
$ heroku pg:backups:url b001
$ heroku pg:backups:download
$ pg_restore -c -d rb3000 latest.dump.1

# 2do

## better logging

starting query run
found x total query results
finishing query run incl x new results

## other

Check if still alive: passive check: if last encounter is older x days (in comparison to last query scan) set active false 

https://www.scraperapi.com/?utm_source=opencollective&utm_medium=github&utm_campaign=colly

better logging: welcher job erzeugt den Eintrag?

Next todo
query=Abrechnung, auf Website 2312 Job, auf mux-search-results nur 250... ?
+ Crawling log mit events


- use "available since" date diretly from job posting  
- lastSeen bei active true check setzen / und nicht mehr bei active false
- while crawling: check first if ID is in DB, then load/open job ad page 
- use codly library for scraping (https://benjamincongdon.me/blog/2018/03/01/Scraping-the-Web-in-Golang-with-Colly-and-Goquery/)

- analytics of duplicates?

- backup DB https://devcenter.heroku.com/articles/heroku-postgres-backups

in tests:
- assert that job names without results have the correct negative string
- assert that healthy job lists always have 26 entries

implement throttel detection
- when multiple monster get requests fail in a row stop and send an alert

analyse
- duplicates
- how many new per day / query
- how many closed per day / query?

## Scan Loop

- read complete ad list (paginated)
- check for sources IDs in URLs/entries
- check if IDs are in DB
- if in DB

## Migrations

### 01 added Stepstone

00: stop everything + BACKUPDB!!

01: psql: ALTER TABLE monster_job_ad_models RENAME COLUMN monster_job_id TO job_source_id;

02: Automigrate (creates job_source)

03: psql: UPDATE monster_job_ad_models SET job_source = 'monster';


### 00 Interface

- getJobAdListForQuery(query string) []string // urls, log Result {query,time, result count}
- getNewAnOldJobsFromJobAdList(jobAdList []string) []string, []string // extract ID, check if ID is in DB, sort into two lists
- markJobAdDBEntryAsSeen(url string) // for all old jobs
- saveJobAdDetails(url string) // for all new jobs
- //checkIfJobAdIsAlive(id string)
- 
