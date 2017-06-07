import Ember from 'ember';

const {
  Helper
} = Ember;

export function date(params/*, hash*/) {
  let createdAt = params[0].split('T');
  createdAt[1] = createdAt[1].substring(0, 8);
  if (params[1]) {
    return createdAt.join(' ');
  }
  return createdAt[0];
}

export default Helper.helper(date);
