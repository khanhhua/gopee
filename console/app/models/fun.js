import DS from 'ember-data';

export default DS.Model.extend({
  fnName: DS.attr('string'),
  xlsxFile: DS.attr('string'),
  inputParams: DS.attr(),
  outputParams: DS.attr()
});
