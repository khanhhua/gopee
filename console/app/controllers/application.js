import Controller from '@ember/controller';

export default Controller.extend({
  isLogin: Ember.computed(function() {
    return !!window.localStorage.getItem('accesstoken');
  }),
  actions: {
    login() {

    },
    closeLoginModal() {

    }
  }
});
