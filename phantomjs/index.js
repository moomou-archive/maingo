var system = require('system');
var server = require('webserver').create();

var INTERVAL_MS = 500;
var ENV = system.env.ENV || 'DEBUG';
var PORT = system.env.PORT || 3333;

var DEBUG = ENV === 'DEBUG';

function pageStyle() {
  return [
    '<style>',
    '#container { display: flex; justify-content: center; }',
    '#canvas { display: -webkit-inline-flex; -webkit-justify-content: center; -webkit-align-items: center;}',
    '#canvas img { float: left; margin: 10px; }',
    '#watermark { position: absolute; margin-left: auto; margin-right: auto; left: 0; right: 0; width: 100px; text-align: center; font-size: 12px; }',
    '</style>'
  ].join('\n');
}

function imgStyle(input) {
  return [
    'display: block;',
    'max-width: ' + (input.width || 100) + 'px;',
    'width: auto;',
    'height: auto;'
  ].join('');
}

function renderPageContent(input) {
  var imgTags = input.urls.map(function imgTag(url) {
    return '<img src="' + url + '" style="' + imgStyle(input) + '"/>';
  }).join('');

  return [
    '<html><head>',
    pageStyle(),
    '</head><body>',
    '<div id="container">',
      '<div id="canvas">',
      imgTags,
      '</div>',
    '</div>',
    '</body></html>'
  ].join('\n');
}

console.log('Listening on', PORT);
server.listen(PORT, function(request, response) {
  var page = new WebPage();
  var body = JSON.parse(request.post);
  var watermark = body.watermark;
  var intervalIdx = null;

  if (DEBUG) console.log('Processing: ', request.post);

  page.viewportSize = { width: 1920, height: 1080 };

  page.onResourceError = function(resourceError) {
    page.error = {
      msg: resourceError.errorString,
      url: resourceError.url,
    };
    console.log('PJ error @ ', page.error.url, page.error.msg);
  };

  page.evaluate(function() {
    document.body.bgColor = 'white'; // 'transparent';
  });

  intervalIdx = window.setInterval(function() {
    var loadingComplete = page.evaluate(function() {
      return window.imagesNotLoaded === 0;
    });

    if (!loadingComplete) return;
    // clear interval immediately to prevent duplication
    clearInterval(intervalIdx);

    // Enable retina screenshot
    page.evaluate(function() {
      /* scale the whole body */
      document.body.style.webkitTransform = 'scale(2)';
      document.body.style.webkitTransformOrigin = '0% 0%';
      /* fix the body width that overflows out of the viewport */
      document.body.style.width = '50%';
    });

    var boundingRect = page.evaluate(function() {
      return document.getElementById('canvas').getBoundingClientRect();
    });

    if (DEBUG) console.log(JSON.stringify(boundingRect));
    page.clipRect = boundingRect;

    if (DEBUG) page.render('output.png');
    var base64 = page.renderBase64('PNG');

    response.statusCode = 200;
    response.write(base64);
    response.close();
    page.close();
  }, INTERVAL_MS);

  var pageContent = renderPageContent(body);
  if (DEBUG) {
    console.log(pageContent);
    console.log('\n');
  }
  page.content = pageContent;

  page.evaluate(function() {
    var images = document.getElementsByTagName('img');
    images = Array.prototype.filter.call(images, function(i) { return !i.complete; });
    window.imagesNotLoaded = images.length;
    Array.prototype.forEach.call(images, function(i) {
      i.onload = function() { window.imagesNotLoaded--; };
    });
  });
});
