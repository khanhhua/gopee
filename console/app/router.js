import EmberRouter from '@ember/routing/router';
import config from './config/environment';

const Router = EmberRouter.extend({
  location: config.locationType,
  rootURL: config.rootURL
});

Router.map(function() {
  this.route('console', function() {
    this.route('edit', { path: 'compose/:id' });
    this.route('compose');
  });
});

export default Router;
