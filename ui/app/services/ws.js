import Ember from 'ember';
import ENV from './../config/environment';

const {
	Service,
	inject,
	set,
	get,
	run
} = Ember;

const RECONECT_DELAY = 3000;

export default Service.extend({
	websockets: inject.service(),

	_ws: null,
  isReady: false,
  msgQueue: [],
	currentChannel: '',

	setupWSConn() {
	let ws = get(this, 'websockets').socketFor(ENV.WSAPIURL);

		ws.on('open', this.openHandler, this);
		ws.on('message', this.messageHandler, this);
		ws.on('close', this.closeHandler, this);
		set(this, '_ws', ws);
	},

	openHandler(event) {
    let msgQueue = get(this, 'msgQueue');
	  set(this, 'isReady', true);
    if (msgQueue.length) {
      for (let msg of msgQueue) {
        get(this, '_ws').send(...msg);
      }
    }
    set(this, 'msgQueue', []);
		if (this.hasOwnProperty('_openHandler')) {
			return this._openHandler(event);
		}
	},
	messageHandler(event) {
		if (this.hasOwnProperty('_messageHandler')) {
		  return this._messageHandler.call(get(this, 'context'), event);
    }
	},
	closeHandler(event) {
    set(this, 'isReady', false);
		if (this.hasOwnProperty('_closeHandler')) {
			return this._closeHandler(event);
		}
		run.later(this, () => {
			get(this, '_ws').reconnect();
		}, RECONECT_DELAY);
	},

	onOpen(callback) {
		set(this, '_openHandler', callback);
	},
	onMessage(callback, context) {
		if (!context) {
			console.log('deprecated! you should pass the context');
			context = this;	// backward compatible
		}
		set(this, '_messageHandler', callback);
		set(this, 'context', context);
	},
	onClose(callback) {
		set(this, '_closeHandler', callback);
	},

	send(...args) {
	  if (!get(this, 'isReady')) {
      get(this, 'msgQueue').pushObject(args);
      return;
    }
    return get(this, '_ws').send(...args);
	},

	subscribeTo(channelName) {
		let currentChannel = get(this, 'currentChannel');
		if (currentChannel.length && currentChannel===channelName) {
			return;
		}
		if (currentChannel.length) {
			this.unsubscribeFrom(currentChannel);
		}

		let action = {
			'type': 'subscribe.to',
			'data': channelName
		};
		this.send({ action }, true);
		set(this, 'currentChannel', channelName);
	},

	unsubscribeFrom(channelName) {
		let action = {
			'type': 'unsubscribe.from',
			'data': channelName
		};
		this.send({ action }, true);
	}
});
