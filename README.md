# Social Tournament Service

## Prerequisites:
* docker & docker-compose

## Run locally

### Clone
*************************************
Clone project from github into local directory using command: `git clone https://github.com/a-andrew/social-tournament-service.git`

OR

Pull build from hub.docker.com using command: `docker pull aandrew3/social-tournament-service`
along with `docker-compose.yml` file and put sources into local directory
*************************************

### Up
Go to folder with the service and run service using `docker-compose up`
Then site will be run on host [http://127.0.0.1]

## Endpoints description
* `/take` - deny:
    * taking all points from player's balance (_no player balance ever goes to zero_)

* `/fund` - deny:
    * creating player with zero balance (_no player balance ever goes to zero_)

* `/announceTournament` - deny:
    * announcing one tournament twice
    
* `/joinTournament` - deny:
    * joining into already finished tournament
    * joining into one tournament twice or being backer twice in one tournament
    * being and player and backer in one tournament
    * joining into tournament if player's balance is less that entry fee
    
* `/resultTournament` - deny:
    * sending result for one tournament twice
    * sending result for already finished tournament
    
* `/balance` - returns player's balance

* `/reset` - reset DB

## Note

* there is controllable environment variable called `_MIN_PLAYER_POINTS_AMOUNT_` in `docker-compose.yml`
  now it set to 100 points but default is 50
  this var specify minimum points player can have
  e.g. MIN_PLAYER_POINTS_AMOUNT = 50, entry fee = 300, player's points = 300. So player have enough points to join into tournament,
  but he shouldn't have zero points balance. So system gives him 50 points compensation.
  Or if player points = 333, then after paying a fee player's balance also will equal 50 (In order to keep the balance no lower than the minimum limit)
  This var can be overridden, and even set to zero
  
If there are some remarks please let me know and I'll fix the logic asap.