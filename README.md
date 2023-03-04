# assignment1

## Introduction 

Combines information about university and country from two different apis. <br>
The apis used to retrieve information are listed below:
<ul>
  <li>University api: http://universities.hipolabs.com/</li>  
  <li>Country api: https://restcountries.com/</li>  
</ul>

## Instructions
### Search for university
&nbsp;&nbsp;&nbsp;&nbsp; 
Path: /unisearcher/v1/uniinfo/{:partial_or_complete_university_name}/
&nbsp;&nbsp;&nbsp;&nbsp;
#### &nbsp; Example:
&nbsp;&nbsp;&nbsp;&nbsp;
/unisearcher/v1/uniinfo/Norwegian University of Science and Technology/
### Search for universities in bordering countries
&nbsp;&nbsp;&nbsp;&nbsp; 
Path: /unisearcher/v1/neighbourunis/{:country_name}/{:partial_or_complete_university_name}{?limit={:number}}/
&nbsp;&nbsp;&nbsp;&nbsp;
#### &nbsp; Example:
&nbsp;&nbsp;&nbsp;&nbsp;
/unisearcher/v1/neighbourunis/United States/University of?limit=5/
### Diagnostics
&nbsp;&nbsp;&nbsp;&nbsp;
Path: /unisearcher/v1/diag/<br>
&nbsp;&nbsp;&nbsp;&nbsp; Gets information about the apis that are used, the version of this service, as well as time since last service restart.

## Decisions 
&nbsp;&nbsp; When retrieving universities from http://universities.hipolabs.com/, it is not an option to search by alpha_two_code. <br>
&nbsp;&nbsp; The country name in http://universities.hipolabs.com/, and the country name in https://restcountries.com/, are not always similar.<br>
&nbsp;&nbsp; Therefore, it is not possible to retrieve both universities and country based on country name.<br>
&nbsp;&nbsp; For that reason, when retrieving universities, all universities are retrieved from the university api based on the university name
&nbsp;&nbsp; provided by the user.<br> Then in the code, universities are filtered based on the country's alpha_two_code.<br><br>
&nbsp;&nbsp; Another way to solve this problem is retrieve all the countries, and retrieve all the universities. Where the isocode matches<br>
&nbsp;&nbsp; but the country name is different, add the name of the country in the country api in a map as the key, and add the country name in<br>
&nbsp;&nbsp; the university api as value. This way, a mapping between a country with different names in the two apis has been created.
