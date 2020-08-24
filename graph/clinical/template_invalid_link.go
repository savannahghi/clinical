package clinical

const invalidLinkTemplate = `
<!doctype html>
<html lang="en">
  <head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css" integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
    <link rel="stylesheet" href="base.css">
    <link rel="stylesheet" href="invalid.css">

    <title>Invalid link</title>
  </head>

  <body>

	<div id="notfound">
		<div class="notfound">
			<div class="notfound-404">
				<h1>Oops!</h1>
				<h2>404 - The Page can't be found</h2>
			</div>
			<!-- <a href="#">Go TO Homepage</a> -->
		</div>
	</div>
	<style>
		* {
  -webkit-box-sizing: border-box;
          box-sizing: border-box;
      }

      body {
        padding: 0;
        margin: 0;
      }

      #notfound {
        position: relative;
        height: 100vh;
      }

      #notfound .notfound {
        position: absolute;
        left: 50%;
        top: 50%;
        -webkit-transform: translate(-50%, -50%);
            -ms-transform: translate(-50%, -50%);
                transform: translate(-50%, -50%);
      }

      .notfound {
        max-width: 520px;
        width: 100%;
        line-height: 1.4;
        text-align: center;
      }

      .notfound .notfound-404 {
        position: relative;
        height: 200px;
        margin: 0px auto 20px;
        z-index: -1;
      }

      .notfound .notfound-404 h1 {
        font-family: 'Montserrat', sans-serif;
        font-size: 236px;
        font-weight: 200;
        margin: 0px;
        color: #211b19;
        text-transform: uppercase;
        position: absolute;
        left: 50%;
        top: 50%;
        -webkit-transform: translate(-50%, -50%);
            -ms-transform: translate(-50%, -50%);
                transform: translate(-50%, -50%);
      }

      .notfound .notfound-404 h2 {
        font-family: 'Montserrat', sans-serif;
        font-size: 28px;
        font-weight: 400;
        text-transform: uppercase;
        color: #211b19;
        background: #fff;
        padding: 10px 5px;
        margin: auto;
        display: inline-block;
        position: absolute;
        bottom: 0px;
        left: 0;
        right: 0;
      }

      .notfound a {
        font-family: 'Montserrat', sans-serif;
        display: inline-block;
        font-weight: 700;
        text-decoration: none;
        color: #fff;
        text-transform: uppercase;
        padding: 13px 23px;
        background: #ff6300;
        font-size: 18px;
        -webkit-transition: 0.2s all;
        transition: 0.2s all;
      }

      .notfound a:hover {
        color: #ff6300;
        background: #211b19;
      }

      @media only screen and (max-width: 767px) {
        .notfound .notfound-404 h1 {
          font-size: 148px;
        }
      }

      @media only screen and (max-width: 480px) {
        .notfound .notfound-404 {
          height: 148px;
          margin: 0px auto 10px;
        }
        .notfound .notfound-404 h1 {
          font-size: 86px;
        }
        .notfound .notfound-404 h2 {
          font-size: 16px;
        }
        .notfound a {
          padding: 7px 15px;
          font-size: 14px;
        }
      }

	</style>

</body>

</html>
`
