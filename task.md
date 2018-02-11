Task Description
================

Overview
--------
Write a web.py ( http://webpy.org) app that exposes a RESTful API for
controlling and retrieving the state of the thermostats in a home.  In a
real-world app you would be issuing commands to real thermostats.  For
this app, you can simply set and store the thermostat's state within the
app and return that state when requested.  Everything can be done in
memory (there is no need for persistence) and you can assume access from
a single client (there is no need to deal with concurrency problems).

Thermostats have the following properties:
------------------------------------------
 ID (read only): unique system identifier for this thermostat (i.e 100
 or 101)

 Name (read-write): display name (i.e "Upstairs Thermostat" and
 "Downstairs Thermostat")

 Current Temp (read-only): since this is not a real home with real
 thermostats an appropriate random value can be returned for this
 property

 Operating Mode (read-write): one of "cool", "heat", "off"

 Cool SetPoint (read-write): a value between 30-100 degrees fahrenheit

 Heat SetPoint (read-write): a value between 30-100 degrees fahrenheit

 Fan Mode (read-write): either "off" or "auto"

Details
-----------
You can assume this app is only controlling the thermostats of a single
home and that home has 2 thermostats.  You need to provide an API to
list all the thermostats in the home and their current state (you can
hard-code appropriate defaults at app startup).  You need to provide
APIs to edit the properties of each thermostat and query the state of
those properties individually.  Appropriate errors should get returned
for accessing non-existent thermostats and non-existent properties.

I am being intentionally vague on what the specifics of the API should
look like.  I want to know how you think the API should look and work.
I will be testing it with curl or something similar.

If you are not familiar with web.py that is fine.  I chose web.py
because it is easy to use and the the tutorial at webpy.org should help
you get going fairly quickly.

We will create a private Slack channel for this assignment.  Don't
hesitate to ask questions.  Take the time you need to complete this, but
I hope this doesn't take an inordinate amount of time.  Please submit
your solution through slack.
