App.Recipe = DS.Model.extend({
    name: DS.attr('string'),
    meal: DS.attr('string'),
    cuisine: DS.attr('string'),
    primary: DS.attr('string'),
    preptime: DS.attr('number'),
    cooktime: DS.attr('number'),
    servings: DS.attr('number'),
});