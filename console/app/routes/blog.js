import Route from '@ember/routing/route';

export default Route.extend({
  model() {
    return Ember.$.ajax(`/call/getArticle`, {
      method: 'POST',
      data: JSON.stringify({
        slug: 'why-run-excel-formula'
      }),
      headers: {
        'x-client-key': '91931784'
      },
      contentType: 'application/json',
      dataType: 'json'
    });
  }
});
