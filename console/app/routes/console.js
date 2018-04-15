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
      }).promise().then(UID => {
        window.localStorage.setItem('x-client-key', UID)
        return true;
      }).catch(err => {
        transition.abort();
      });
    } else {
      return true;
    }
  }
});
