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
When retrieving universities from http://universities.hipolabs.com/, it is not an option to search by alpha_two_code. <br>
The country name in http://universities.hipolabs.com/, and the country name in https://restcountries.com/, are not always similar.<br>
Therefore, it is not possible to retrieve both universities and country based on country name.<br>
For that reason, when retrieving universities, all universities are retrieved from the university api based on the university name
provided by the user. Then in the code, universities are filtered based on the country's alpha_two_code.<br><br>

Another way to solve this problem is to retrieve all the countries, and retrieve all the universities. Where the isocode matches<br>
but where the country name is different, add the country name from the country api in a map as the key, and add the country name from<br>
the university api as value. This way, a mapping between a country with different names in the two apis has been created.<br><br>

The reasoning for why the first solution was chosen, is that the second solution can be vulnerable to changes in the apis.<br>
If one of the apis changes a country name, and the map in this service has not yet been updated, the mapping will not work anymore.<br>
The first solution can therefore be more robust in terms of changes in the apis that are used.<br><br>

The first solution is still not perfect, as a lot of unnecessary data can be retrieved as all universities with a said name is retrieved,<br>
instead of just the universities with a specific name in specific countries. Alt