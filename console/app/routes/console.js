import Route from '@ember/routing/route';
import Ember from "ember";

export default Route.extend({
  beforeModel(transition) {
    const { code } = transition.queryParams;
    if (code) {
      return Ember.$.ajax('/api/token', {
        method: 'POST',
        contentType: 'application/json',
        data: JSON.stringify({
          code
        })
      }).promise().then(token => {
        window.localStorage.setItem('accesstoken', token)
        window.location.href = window.location.origin + '/console';
      }).catch(err => {
        transition.abort();
      });
    } else {
      return true;
    }
  }
});
