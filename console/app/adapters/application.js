import DS from 'ember-data';

export default DS.RESTAdapter.extend({
  namespace: 'api',
  headers: Ember.computed(function () {
    if (!!localStorage.getItem('accesstoken')) {
      return {
        'Authorization': `Bearer ${localStorage.getItem('accesstoken')}`
      }
    } else {
      return {};
    }
  })
});
