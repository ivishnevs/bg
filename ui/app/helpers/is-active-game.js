import Ember from 'ember';

const {
  Helper
} = Ember;

export function isActiveGame(params/*, hash*/) {
  let game = params[0];
  if (game.status !== 'closed' && game.gamerCount !== game.occupiedPlaces) {
    return true;
  }
}

export default Helper.helper(isActiveGame);
