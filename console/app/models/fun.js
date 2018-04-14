import DS from 'ember-data';

export default DS.Model.extend({
  fnName: DS.attr('string'),
  xlsxName: DS.attr('string'),
  inputMappings: DS.attr(),
  outputMappings: DS.attr()
});
