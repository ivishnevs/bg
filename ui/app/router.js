import Ember from 'ember';
import config from './config/environment';

const Router = Ember.Router.extend({
  location: config.locationType,
  rootURL: config.rootURL
});

Router.map(function() {
  this.route('sign-in');
  this.route('sign-up');
  this.route('rooms', function() {
    this.route('details', {path: '/:room_id'}, function() {
      this.route('gameroles', {path: '/:game_id'});
    });
  });
  this.route('start-page');
  this.route('about');
  this.route('admin');
  this.route('gamer', {path: '/gamer/:gamer_id'});
  this.route('statistics', {path: '/statistics/:game_id'});
});

export default Router;
