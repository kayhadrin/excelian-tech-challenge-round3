@echo off
setlocal
set ITERATIONS=10

echo ITERATIONS=%ITERATIONS%

set STARTTIME=%time%

for /l %%x in (1, 1, %ITERATIONS%) do (
	bin\multiCoreMin.exe > nul
)

set ENDTIME=%time%

echo Start: %STARTTIME%
echo End: %ENDTIME%