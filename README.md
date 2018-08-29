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

# hot reload

$ gin

# heroku pg backups
$ heroku pg:backups:capture
$ heroku pg:backups:url b001
$ heroku pg:backups:download
$ pg_restore -c -d rb3000 latest.dump.1

# 2do

in tests:
- assert that job names without results have the correct negative string
- assert that healthy job lists always have 26 entries

implement throttel detection
- when multiple monster get requests fail in a row stop and send an alert