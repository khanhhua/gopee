import Route from '@ember/routing/route';

export default Route.extend({
  actions: {
    loginWithDropbox() {
      window.open('https://www.dropbox.com/oauth2/authorize?' +
        'client_id=j4365xi2ynl3zri'+
        '&response_type=code'+
        '&redirect_uri=' + window.location.origin + '/console', '_blank');
    },
    logout() {
      window.localStorage.removeItem('accesstoken');
      window.location.href = window.location.origin + '/'
    }
  }
});
