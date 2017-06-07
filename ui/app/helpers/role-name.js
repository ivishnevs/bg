import Ember from 'ember';

const {
	Helper
} = Ember;

export function roleName(params/*, hash*/) {
  let role = params[0];
  let gamerCount = params[1];

  if (gamerCount <= 4) {
		switch (role) {
			case 0:
				return "Розничный продавец";
			case 1:
				return "Оптовый продавец";
			case 2:
				return "Дистрибьютор";
			case 3:
				return "Завод";
		}
  }
  switch (gamerCount - role) {
		case gamerCount:
			return "Розничный продавец";
		case 2:
			return "Дистрибьютор";
		case 1:
			return "Завод";
    default:
      return `Оптовый продавец #${role}`;
  }
}

export default Helper.helper(roleName);
