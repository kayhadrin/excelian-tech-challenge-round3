@echo off
setlocal
set ITERATIONS=20

echo ITERATIONS=%ITERATIONS%

set STARTTIME=%time%

for /l %%x in (1, 1, %ITERATIONS%) do (
	bin\multiCoreXor.exe > nul
)

set ENDTIME=%time%

echo Start: %STARTTIME%
echo End: %ENDTIME%
