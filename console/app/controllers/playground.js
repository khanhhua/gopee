import Controller from '@ember/controller';
import { computed } from '@ember/object';


export default Controller.extend({
  isLoading: false,
  actualInputParams: computed('fun.inputParams', function () {
    if (!this.get('fun.inputParams.[]')) {
      return [];
    }

    return this.get('fun.inputParams.[]').map(item => ({
      name: item.name,
      address: item.address,
      value: ''
    }));
  }),
  actions: {
    selectFun(fun) {
      this.setProperties({
        fun,
        isLoading: true
      });

      this.get('store').find('fun', fun.id).then(fun => {
        this.setProperties({
          fun,
          isLoading: false
        });
      });
    },
    callFunction() {
      this.set('isLoading', true);

      const data = this.get('actualInputParams.[]').reduce((acc, item) => {
        acc[item.name] = item.value;
        return acc;
      }, {});
      let clientKey;
      try {
        clientKey = JSON.parse(atob(window.localStorage.getItem('accesstoken').split('.')[1])).sub;
        Ember.$.ajax(`/call/${this.get('fun.fnName')}`, {
          headers: {
            'x-client-key': clientKey
          },
          data: JSON.stringify(data),
          method: 'POST',
          contentType: 'application/json'
        }).promise().then(response => {
          this.set('funOutput', JSON.stringify(response, 4, true))
          this.set('isLoading', false);
        }).catch(() => {
          this.set('isLoading', false);
        })
      } catch (e) {
        console.warn('Authentication error')
      }
    }
  }
});
