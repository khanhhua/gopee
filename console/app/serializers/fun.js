import DS from 'ember-data';

export default DS.RESTSerializer.extend(DS.EmbeddedRecordsMixin, {
  attrs: {
    inputParams: { embedded: { serialize: true, deserialize: true }},
    outputParams: { embedded: { serialize: true, deserialize: true }}
  }
});
