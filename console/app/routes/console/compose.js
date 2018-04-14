import Route from '@ember/routing/route';

export default Route.extend({
  model() {
    return {
      fnName: '',
      xlsxFile: '',
      inputParams: [
        {
          name: '',
          address: ''
        }
      ],
      outputParams: [
        {
          name: '',
          address: ''
        }
      ]
    };
  },
  actions: {
    save(model) {
      console.log(`Saving model:`, model)
      const fun = this.get('store').createRecord('fun', model);
      fun.save().catch(err => {
        fun.unloadRecord()
      });
    },
    addInputMapping() {
      this.currentModel.inputParams.pushObject({
        name: '',
        address: ''
      });
    },
    removeInputMapping(nth) {
      this.currentModel.inputParams.removeAt(nth);
    },
    addOutputMapping() {
      this.currentModel.outputParams.pushObject({
        name: '',
        address: ''
      });
    },
    removeOutputMapping(nth) {
      this.currentModel.outputParams.removeAt(nth);
    },
  }
});
