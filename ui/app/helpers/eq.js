import Ember from 'ember';

const {
  Helper
} = Ember;

export function increment(params/*, hash*/) {
  return params[0] === params[1];
}

export default Helper.helper(increment);
