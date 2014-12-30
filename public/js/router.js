App.Router.map(function() {
	this.resource('recipes', function() {
		this.resource('recipe', { path: ':recipe_id' });
	});
});

App.IndexRoute = Ember.Route.extend({
	beforeModel: function() {
		this.transitionTo('recipes');
	}
});

App.RecipesRoute = Ember.Route.extend({
	model: function() {
		return this.store.find('recipe');
	}
});

App.RecipeRoute = Ember.Route.extend({
	model: function(params) {
		return this.store.find('recipe', params.recipe_id);
	}
})