import Ember from 'ember';

const {
  Helper
} = Ember;

export function isCurrentGamer(params/*, hash*/) {
  return params[0] === params[1];
}

export default Helper.helper(isCurrentGamer);
