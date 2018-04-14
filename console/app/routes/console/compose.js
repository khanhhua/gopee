import Route from '@ember/routing/route';

export default Route.extend({
  model() {
    return {
      fnName: '',
      xlsxFile: '',
      inputMappings: [
        {
          name: '',
          address: ''
        }
      ],
      outputMappings: [
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
      fun.save();
    },
    addInputMapping() {
      this.currentModel.inputMappings.pushObject({
        name: '',
        address: ''
      });
    },
    removeInputMapping(nth) {
      this.currentModel.inputMappings.removeAt(nth);
    },
    addOutputMapping() {
      this.currentModel.outputMappings.pushObject({
        name: '',
        address: ''
      });
    },
    removeOutputMapping(nth) {
      this.currentModel.outputMappings.removeAt(nth);
    },
  }
});
