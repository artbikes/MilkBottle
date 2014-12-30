(function() {
	
    App = Ember.Application.create({
    	  LOG_TRANSITIONS: true 	
    });

    App.ApplicationStore = DS.Store.extend({
      revision: 12,
       adapter: 'App.Adapter'

    });
   
    App.ApplicationAdapter = DS.RESTAdapter.extend({
        namespace: 'api'
    });

})();
