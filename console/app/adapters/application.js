import DS from 'ember-data';

export default DS.RESTAdapter.extend({
  namespace: 'api',
  headers: Ember.computed(function () {
    return {
      'x-client-key': localStorage.getItem('x-client-key')
    }
  })
});
