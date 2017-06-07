import Ember from 'ember';

const {
  Helper
} =Ember;

export function getLength(params/*, hash*/) {
  return params[0].length - parseInt(params[1]);
}

export default Helper.helper(getLength);
