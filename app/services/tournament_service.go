package services

import (
	"github.com/social-tournament-service/app/daos"
	"github.com/social-tournament-service/app/models"
	"errors"
	"fmt"
	"gopkg.in/pg.v5"
)

type TournamentService struct{
	playerDao *daos.PlayerDao
	tournamentDao *daos.TournamentDao
}

func NewTournamentService(tournamentDao *daos.TournamentDao, playerDao *daos.PlayerDao) *TournamentService{
	return &TournamentService{
		playerDao: playerDao,
		tournamentDao: tournamentDao,
	}
}

func (s *TournamentService) Announce(tournamentId string, deposit int) error{
	return s.tournamentDao.Create(&models.Tournament{
		Id: tournamentId,
		Deposit: deposit,
	})
}

func (s *TournamentService) Join(tournamentId, playerId string, backerIds []string) error{
	backerIds = getUniqueMemberIds(backerIds)
	
	if playerIsBacker(playerId, backerIds){
		return errors.New("Player can't be player and backer in one tournament")
	}
	
	memberIds := append(backerIds, playerId)
	if err := s.validateForJoining(tournamentId, memberIds); err != nil{
		return err
	}
	
	entryFee, err := s.calculateEntryFee(tournamentId, memberIds)
	if err != nil{
		return err
	}

	members, err := s.getMembersBalance(memberIds)
	if err := s.checkIfPlayersCanPayEntryFee(entryFee, tournamentId, members); err != nil{
		return err
	}

	return s.joinToTournament(tournamentId, playerId, backerIds, members, entryFee)
}

func (s *TournamentService) Result(result *models.TournamentResultJson) error{
	if err := s.validateForRewarding(result); err != nil{
		return err
	}
	
	return s.finishTournament(result)
}

func (s *TournamentService) validateForJoining(tournamentId string, memberIds []string) error{
	//does tournament exist
	if !s.tournamentDao.ExistsAndNotFinished(&models.Tournament{Id: tournamentId}){
		return errors.New(fmt.Sprintf("Tournament with ID `%s` not found or is already finished", tournamentId))
	}
	
	//validate members
	for _, memberId := range memberIds{
		//does member exist
		if !s.playerDao.Exists(&models.Player{Id: memberId}){
			return errors.New(fmt.Sprintf("Member with ID `%s` not found", memberId))
		}

		tBacker := &models.TournamentPlayer{
			PlayerId: memberId,
			TournamentId: tournamentId,
		}
		
		//has member already joined to the tournament as a player
		if s.tournamentDao.IsJoinedAsPlayer(tBacker){
			return errors.New(fmt.Sprintf("Member with ID `%s` already has joided as a player to tournament with ID `%s`", memberId, tournamentId))
		}

		//has member already joined to the tournament as a backer
		if s.tournamentDao.IsJoinedAsBacker(tBacker){
			return errors.New(fmt.Sprintf("Member with ID `%s` already has joided as a backer to tournament with ID `%s`", memberId, tournamentId))
		}
	}
	
	return nil
}

func (s *TournamentService) validateForRewarding(result *models.TournamentResultJson) error{
	//does tournament exist
	if !s.tournamentDao.ExistsAndNotFinished(&models.Tournament{Id: result.TournamentId}){
		return errors.New(fmt.Sprintf("Tournament with ID `%s` not found or is already finished", result.TournamentId))
	}

	//does tournament announced (started)
	if !s.tournamentDao.IsStarted(&models.TournamentPlayer{TournamentId: result.TournamentId}){
		return errors.New(fmt.Sprintf("Tournament with ID `%s` is not started yet", result.TournamentId))
	}

	//does everyone winner exist
	for _, winner := range result.Winners{
		if !s.playerDao.Exists(&models.Player{Id: winner.PlayerId}){
			return errors.New(fmt.Sprintf("Winner with ID `%s` not found", winner.PlayerId))
		}
	}
	
	return nil
}

func (s *TournamentService) finishTournament(result *models.TournamentResultJson) error{
	return s.tournamentDao.Transaction(func(txn *pg.Tx) error {
		if err := s.tournamentDao.FinishTxn(txn, &models.Tournament{
			Id: result.TournamentId,
			Finished: true,
		}); err != nil{
			return err
		}

		return s.rewardWinners(txn, result)
	})
}

func (s *TournamentService) rewardWinners(txn *pg.Tx, result *models.TournamentResultJson) error{
	rewarded := map[string]bool{}
	
	for _, winner := range result.Winners{
		//avoid duplicating players in winners list
		if _, ok := rewarded[winner.PlayerId]; ok{
			continue
		}
		
		if winner.Prize < 0{
			return errors.New("Tournament prize should not be less than zero")
		}

		data := &models.TournamentPlayer{
			TournamentId: result.TournamentId,
			PlayerId: winner.PlayerId,
		}

		if err := s.tournamentDao.GetBackers(data); err != nil{
			return err
		}

		memberIds := append(data.BackerIds, winner.PlayerId)

		membersPrize := 0
		if winner.Prize > 0{
			membersPrize = winner.Prize / len(memberIds)
			if membersPrize == 0{
				membersPrize = 1
			}
		}

		for _, memberId := range memberIds{			
			if err := s.playerDao.AddPointsTxn(txn, &models.Player{
				Id: memberId,
				Points: membersPrize,
			}); err != nil{
				return err
			}
		}

		rewarded[winner.PlayerId] = true
	}
	
	return nil
}

func (s *TournamentService) getMembersBalance(memberIds []string) ([]*models.Player, error){
	members := make([]*models.Player, 0)
	for _, memberId := range memberIds{
		member := &models.Player{Id: memberId}
		if err := s.playerDao.GetBalance(member); err != nil{
			return nil, err
		}
		
		members = append(members, member)
	}
	
	return members, nil
}

func (s *TournamentService) checkIfPlayersCanPayEntryFee(entryFee int, tournamentId string, members []*models.Player) error {
	for _, member := range members{		
		if member.Points < entryFee{
			return errors.New(fmt.Sprintf("Member with ID `%s` doesn't have enough points to pay entry fee", member.Id))
		}
	}
	
	return nil
}

func (s *TournamentService) calculateEntryFee(tournamentId string, membersIds []string) (int, error) {
	tournament := &models.Tournament{Id:tournamentId}
	if err := s.tournamentDao.GetDeposit(tournament); err != nil{
		return 0, err
	}

	if tournament.Deposit == 0{
		return 0, nil
	}
	
	entryFee := int(tournament.Deposit / (len(membersIds)))
	if entryFee == 0{
		entryFee = 1
	}
	
	return entryFee, nil
}

func (s *TournamentService) joinToTournament(tournamentId, playerId string, backerIds []string, members []*models.Player, entryFee int) error{
	return s.tournamentDao.Transaction(func(txn *pg.Tx) error {
		if err := s.tournamentDao.CreateTournamentPlayerTxn(txn, &models.TournamentPlayer{
			PlayerId: playerId,
			TournamentId: tournamentId,
			BackerIds: backerIds,
		}); err != nil{
			return err
		}

		return s.takeEntryFeeFromMembers(txn, entryFee, members)
	})
}

func (s *TournamentService) takeEntryFeeFromMembers(txn *pg.Tx, entryFee int, members []*models.Player) error{
	for _, member := range members{
		member.Points -= entryFee
		if member.Points < daos.MIN_PLAYER_POINTS_AMOUNT{
			member.Points = daos.MIN_PLAYER_POINTS_AMOUNT
		}
		
		if err := s.playerDao.UpdatePointsTxn(txn, member); err != nil{
			return err
		}
	}
	
	return nil
}

func getUniqueMemberIds(ids []string) []string{
	temp := map[string]string{}
	for _, id := range ids{
		temp[id] = id
	}
	
	uniqueIds := []string{}
	for id := range temp{
		uniqueIds = append(uniqueIds, id)
	}
	
	return uniqueIds
}

func playerIsBacker(playerId string, backerIds []string) bool{
	for _, bId := range backerIds{
		if playerId == bId{
			return true
		}
	}
	
	return false
}