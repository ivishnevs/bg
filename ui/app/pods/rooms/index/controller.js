import Ember from 'ember';

const {
  Controller,
  computed,
  get,
  run,
  set
} = Ember;

export default Controller.extend({

  filteredRooms: computed('model', {
    get() {
      return get(this, 'model');
    },
    set(key, value) {
      return value;
    }
  }),

  filterModelByName(searchTerm) {
    set(this, 'filteredRooms', get(this, 'model').filter(
      room => ~room.name.toLowerCase().indexOf(
        searchTerm.toLowerCase()
      )
    ));
  },

  actions: {
    searchInputUpPressHandler() {
      run.debounce(this, 'filterModelByName', get(this, 'searchTerm'), 200);
    },
    roomSelected(roomId) {
      this.transitionToRoute('rooms.details', roomId);
    }
  }
});
