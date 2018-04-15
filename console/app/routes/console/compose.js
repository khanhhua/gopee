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
    browse() {
      const options = {
        // Required. Called when a user selects an item in the Chooser.
        success: (files) => {
          console.debug("Here's the file link:", files[0]);
          Ember.setProperties(this.currentModel, {
            xlsxFile: `/${files[0].name}`, // always relative to app folder
          });
        },

        // Optional. Called when the user closes the dialog without selecting a file
        // and does not include any parameters.
        cancel: function() {

        },

        // Optional. "preview" (default) is a preview link to the document for sharing,
        // "direct" is an expiring link to download the contents of the file. For more
        // information about link types, see Link types below.
        linkType: "preview", // or "direct"

        // Optional. A value of false (default) limits selection to a single file, while
        // true enables multiple file selection.
        multiselect: false, // or true

        // Optional. This is a list of file extensions. If specified, the user will
        // only be able to select files with these extensions. You may also specify
        // file types, such as "video" or "images" in the list. For more information,
        // see File types below. By default, all extensions are allowed.
        extensions: ['.xlsx', '.xls'],

        // Optional. A value of false (default) limits selection to files,
        // while true allows the user to select both folders and files.
        // You cannot specify `linkType: "direct"` when using `folderselect: true`.
        folderselect: false, // or true
      };

      Dropbox.choose(options);
    },
    save(model) {
      const fun = this.get('store').createRecord('fun', model);
      fun.save()
        .then((model) => {
          this.transitionTo('console.edit', model.id)
        })
        .catch(err => {
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
