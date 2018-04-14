import { module, test } from 'qunit';
import { setupTest } from 'ember-qunit';

module('Unit | Route | console/compose', function(hooks) {
  setupTest(hooks);

  test('it exists', function(assert) {
    let route = this.owner.lookup('route:console/compose');
    assert.ok(route);
  });
});
