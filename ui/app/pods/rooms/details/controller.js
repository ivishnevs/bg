import Ember from 'ember';

const {
  Controller,
  computed
} = Ember;

export default Controller.extend({
	gamesSortingDesc: ['CreatedAt:desc'],
	games: computed.sort('model.games', 'gamesSortingDesc'),

  actions: {
    gameSelected(gameId, gameStatus) {
      if (gameStatus === 'closed' || gameStatus === 'in process') {
        return;
      }
      this.transitionToRoute('rooms.details.gameroles', gameId);
    }
  }
});