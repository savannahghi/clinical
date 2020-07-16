package clinical

const patientProfileTemplate = `
<!doctype html>
<html lang="en">

<head>
    <!-- Required meta tags -->
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, initial-scale=1, shrink-to-fit=no">

    <!-- Bootstrap CSS -->
    <link rel="stylesheet" href="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/css/bootstrap.min.css"
				integrity="sha384-Vkoo8x4CGsO3+Hhxv8T/Q5PaXtkKtu6ug5TOeNV6gBiFeWPGFN9MuhOf23Q9Ifjh" crossorigin="anonymous">
		<!-- <link rel="stylesheet" href="base.css">
		<link rel="stylesheet" href="profile.css"> -->
		<!-- Optional JavaScript -->
		<!-- jQuery first, then Popper.js, then Bootstrap JS -->
		<script src="https://code.jquery.com/jquery-3.4.1.slim.min.js"
    integrity="sha384-J6qa4849blE2+poT4WnyKhv5vZF5SrPo0iEjwBvKU7imGFAV0wwj1yYfoRSJoZ+n" crossorigin="anonymous"></script>
		<script src="https://cdn.jsdelivr.net/npm/popper.js@1.16.0/dist/umd/popper.min.js"
     integrity="sha384-Q6E9RHvbIyZFJoft+2mJbHaEWldlvI9IOYy5n3zV9zzTtmI3UksdQRVvoxMfooAo"crossorigin="anonymous"></script>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/4.4.1/js/bootstrap.min.js"
     integrity="sha384-wfSDF2E50Y2D1uUdj0O3uMBJnjuUD4Ih7YwaYd1iqfktj0Uod8GCExl3Og8ifwB6"crossorigin="anonymous"></script>
				
		<!--Google Fonts-->
		<link href="https://fonts.googleapis.com/css2?family=Red+Hat+Display:wght@400;500;700&display=swap" rel="stylesheet">

		<!--FontAwsome-->
		<script src="https://kit.fontawesome.com/a076d05399.js"></script>
		<title>Patient Profile</title>
</head>

<body>
    <section class="header-title">
      <h4>HealthCloud</h4>
		</section>
		<section class="container"> 

			<section class="patient-profile-container">
				<div class="title">
					<h5>Patient profile</h5>
				</div>

				<div class="patient-profile-details-container grid-container">
					<div class="patient-avatar grid-item">
						<i class="fas fa-user-circle fa-7x"></i>
					</div>
					<div class="patient-profile-details grid-item">
						<h6 class="patient-name">{{.Name}}</h6>
            <h6><span><i class="far fa-calendar-alt"></i></span>  {{.Age}}</h6>
            <h6><span><i class="fa fa-venus-mars"></i></span>  {{.Gender}}</h6>
						<h6><span><i class="fas fa-map-marker-alt"></i></span>  {{.Addresses}}</h6>
					</div>
				</div>
			</section>

			<section class="patient-bio-container">
				<div class="title">
					<h5>Patient Bio</h5>
				</div>

				<div class="bio grid-container">
					<div class="grid-item">
						<div class="bio-card allergies">
							<h6>PROBLEMS AND ALLERGIES</h6>
							<h5>{{.Problems}}</h5>
							<h5>{{.Allergies}}</h5>
						</div>
					</div>
					<div class="grid-item">
						<div class="bio-card id-doc">
							<h6>IDENTIFICATION DOCUMENT</h6>
							<h5>{{.IDDocs}}</h5>
						</div>
					</div>
					<div class="grid-item">
						<div class="bio-card languages">
							<h6>LANGUAGES</h6>
							<h5>{{.Languages}}</h5>
						</div>
					</div>
					<div class="grid-item">
						<div class="bio-card marital-status">
							<h6>MARITAL STATUS</h6>
							<h5>{{.MaritalStatus}}</h5>
						</div>
					</div>
				</div>
			</section>
		</section>

        <style>
				* {
  					box-sizing: border-box;
					}

					html {
    				position: relative;
    				min-height: 100%;
					}

          body {
						font-family: 'Red Hat Display', sans-serif;
            background-color: #ffffff;
          }

					.header-title{
						background-color: #603A8B;
  					height: 7vh;
						top: 0;
						position: sticky;
					}

					.header-title h4{
						color: #ffffff;
						padding-top: 20px;
						text-align: center;
					}

					.title h5{
						color: #999999;
					}

					.fa-user-circle{
						color: #999999;
					}

					.patient-profile-container{
						margin-bottom: 40px;
					}

					.patient-bio-container{
						margin-bottom: 70px;
					}
					

					.grid-container{
						display: inline-grid;
						grid-template-columns: auto auto;
						grid-column-gap: 50px;
						grid-row-gap: 50px;
					}

					.patient-name{
						margin-top: 5px;
					}

					.patient-profile-details h6{
						padding-bottom: 10px;
					}

					.bio-card{
						border-radius: 15px;
						background-color: #F6F7F9;
						padding: 20px 30px;
						width: 520px;
					}

					.bio-card h6{
						color: #999999;
						padding-bottom: 15px;
					}
					
					.bio-card h5{
						color: #603A8B;
					}

					.allergies{
						background-color: #E0EEE8;
					}

					.id-doc{
						background-color: #F6F7F9;
					}

					.languages{
						background-color: #FAF0EA;
					}

					.marital-status{
						background-color: #E3F0FF;
					}
					
          .title {
            font-weight: normal;
            padding-top: 30px;
            padding-bottom: 30px;
          }

					footer {
    				background-color: #f2f2f2;
    				position: absolute;
    				left: 0;
    				bottom: 0;
    				height: 60px;
    				width: 100%;
    				overflow: hidden;
						text-align: center;
						padding-top: 15px;
					}
					

					@media only screen and (max-width: 920px) {
						.grid-container{
							grid-gap: 20px;
  						grid-template-columns: repeat(auto-fit, 186px);
						}

						.bio-card{
							padding: 20px 30px;
							width: 380px;
						}
					}
				</style>
				<footer>
					<small>
						© 2020 Savannah Informatics - All Rights Reserved 
					</small>
				</footer>
    </body>
</html>
`
