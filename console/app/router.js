import EmberRouter from '@ember/routing/router';
import config from './config/environment';

const Router = EmberRouter.extend({
  location: config.locationType,
  rootURL: config.rootURL
});

Router.map(function() {
  this.route('console', function() {
    this.route('index', { path: '/' });
    this.route('edit', { path: 'compose/:id' });
    this.route('compose');
  });
  this.route('playground', { path: 'console/playground' });
  this.route('blog', { path: '/' });
});

export default Router;
