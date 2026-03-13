package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/Polshkrev/gopolutils"
	"github.com/Polshkrev/gopolutils/fayl"
)

const (
	logs       string = "./logs"
	dateFormat string = "2006-01-02"
)

func emptyString() string {
	return ""
}

func isEmptyString(value string) bool {
	return len(value) == 0 || value == emptyString()
}

func trimExtension(path string) (string, *gopolutils.Exception) {
	var cleaned string = filepath.Clean(path)
	var extension string = filepath.Ext(cleaned)
	var stripped string = strings.TrimSuffix(cleaned, extension)
	if len(stripped) == len(cleaned) {
		return path, gopolutils.NewNamedException(gopolutils.ValueError, "Can not strip '%s' from '%s'.", extension, cleaned)
	}
	return stripped, nil
}

func splitFileName(path *fayl.Path) (*fayl.Path, *gopolutils.Exception) {
	var pathString string = path.ToString()
	var suffix fayl.Suffix
	var suffixException *gopolutils.Exception
	suffix, suffixException = path.Suffix()
	if suffixException != nil {
		return nil, suffixException
	}
	var stripped string
	var strippedException *gopolutils.Exception
	stripped, strippedException = trimExtension(pathString)
	if strippedException != nil {
		return nil, strippedException
	} else if suffix == fayl.Gz {
		var baseName string
		var baseNameException *gopolutils.Exception
		baseName, baseNameException = trimExtension(stripped)
		if baseNameException != nil {
			return nil, baseNameException
		}
		return fayl.PathFrom(baseName), nil
	}
	return fayl.PathFrom(stripped), nil
}

func constructTargetFileName(originalPath *fayl.Path, platform fayl.OS) (string, *gopolutils.Exception) {
	var result string
	var folder string = filepath.Base(originalPath.ToString())
	switch platform {
	case fayl.WINDOWS:
		var suffix string
		var suffixException *gopolutils.Exception
		suffix, suffixException = fayl.StringFromSuffix(fayl.Zip)
		if suffixException != nil {
			return emptyString(), suffixException
		}
		result = fmt.Sprintf("%s.%s", folder, suffix)
	default:
		var tarSuffix string
		var suffixException *gopolutils.Exception
		tarSuffix, suffixException = fayl.StringFromSuffix(fayl.Tar)
		if suffixException != nil {
			return emptyString(), suffixException
		}
		var gzipSuffix string
		var gzipSuffixException *gopolutils.Exception
		gzipSuffix, gzipSuffixException = fayl.StringFromSuffix(fayl.Gz)
		if gzipSuffixException != nil {
			return emptyString(), gzipSuffixException
		}
		result = fmt.Sprintf("%s.%s.%s", folder, tarSuffix, gzipSuffix)
	}
	return result, nil
}

func positionalRangeToFiles(arguments []string, logger *gopolutils.Logger) ([]*fayl.Entry, *gopolutils.Exception) {
	var result []*fayl.Entry
	switch len(arguments) {
	case 0:
		var rootPath *fayl.Path = fayl.NewPath()
		logger.Log(fmt.Sprintf("Loading files from '%s'.", rootPath.ToString()), gopolutils.Info)
		var directory *fayl.Directory = fayl.NewDirectory(rootPath)
		var readExcept *gopolutils.Exception = directory.Read()
		if readExcept != nil {
			return nil, readExcept
		}
		result = directory.Collect()
	case 1:
		var value string = flag.Arg(0)
		var root *fayl.Path = fayl.PathFrom(value)
		logger.Log(fmt.Sprintf("Loading files from '%s'.", root.ToString()), gopolutils.Info)
		var directory *fayl.Directory = fayl.NewDirectory(root)
		var readExcept *gopolutils.Exception = directory.Read()
		if readExcept != nil {
			return nil, readExcept
		}
		result = directory.Collect()
	default:
		var i int
		for i = range arguments {
			var argument string = arguments[i]
			var entry *fayl.Entry = fayl.NewEntry(fayl.PathFrom(argument))
			logger.Log(fmt.Sprintf("Archiving '%s'.", entry.Path().ToString()), gopolutils.Info)
			result = append(result, entry)
		}
	}
	return result, nil
}

func archive(source string, logger *gopolutils.Logger) *gopolutils.Exception {
	if isEmptyString(source) {
		var current *fayl.Path = fayl.NewPath()
		var targetString string
		var targetStringException *gopolutils.Exception
		targetString, targetStringException = constructTargetFileName(current, fayl.OS(runtime.GOOS))
		if targetStringException != nil {
			return targetStringException
		}
		var target *fayl.Path = fayl.PathFrom(targetString)
		logger.Log(fmt.Sprintf("Creating archive '%s'", target.ToString()), gopolutils.Info)
		var files []*fayl.Entry
		var filesException *gopolutils.Exception
		files, filesException = positionalRangeToFiles(flag.Args(), logger)
		if filesException != nil {
			return filesException
		}
		return fayl.ZipFolder(target, files...)
	}
	var destination *fayl.Path = fayl.PathFrom(source)
	logger.Log(fmt.Sprintf("Creating archive '%s'", destination.ToString()), gopolutils.Info)
	var files []*fayl.Entry
	var filesException *gopolutils.Exception
	files, filesException = positionalRangeToFiles(flag.Args(), logger)
	if filesException != nil {
		return filesException
	}
	return fayl.Archive(destination, files...)
}

func createTimeFilename(logFolder string, dateFormat string, suffix fayl.Suffix) (*fayl.Path, *gopolutils.Exception) {
	var suffixString string
	var suffixException *gopolutils.Exception
	suffixString, suffixException = fayl.StringFromSuffix(suffix)
	if suffixException != nil {
		return nil, suffixException
	}
	var now time.Time = time.Now()
	return fayl.PathFrom(fmt.Sprintf("%s%c%s.%s", logFolder, filepath.Separator, now.Format(dateFormat), suffixString)), nil
}

func setupLogger(name, dateFormat string, defaultLevel gopolutils.LoggingLevel, verbose bool) (*gopolutils.Logger, *gopolutils.Exception) {
	var logger *gopolutils.Logger = gopolutils.NewLogger(name, defaultLevel)
	if verbose {
		var except *gopolutils.Exception = logger.AddConsole()
		if except != nil {
			return nil, except
		}
	}
	var filename *fayl.Path
	var filenameException *gopolutils.Exception
	filename, filenameException = createTimeFilename(logs, dateFormat, fayl.Log)
	if filenameException != nil {
		return nil, filenameException
	}
	var parent *fayl.Path
	var parentException *gopolutils.Exception
	parent, parentException = filename.Parent()
	if parentException != nil {
		return nil, parentException
	} else if !parent.Exists() {
		var entry *fayl.Entry = fayl.NewEntry(parent)
		entry.SetType(fayl.DirectoryType)
		var createException *gopolutils.Exception = entry.MakeDirectory()
		if createException != nil {
			return nil, createException
		}
	}
	var fileException *gopolutils.Exception = logger.AddFile(filename.ToString())
	if fileException != nil {
		return nil, fileException
	}
	return logger, nil
}

func main() {
	var create *bool = flag.Bool("c", false, "Create a specified archive on the filesystem.")
	var source *string = flag.String("f", emptyString(), "Specify a source file. This flag is only required if the '-c' flag is not set.")
	var verbose *bool = flag.Bool("v", false, "Display the logs to standard output.")
	flag.Parse()
	var logger *gopolutils.Logger
	var loggerException *gopolutils.Exception
	logger, loggerException = setupLogger("ar", dateFormat, gopolutils.Debug, *verbose)
	if loggerException != nil {
		panic(loggerException)
	} else if isEmptyString(*source) && !(*create) {
		if *verbose {
			logger.Log(gopolutils.NewNamedException(gopolutils.RuntimeError, "No source file(s) have been provided.").Error(), gopolutils.Critical)
			flag.Usage()
		}
		os.Exit(1)
	} else if *create {
		var except *gopolutils.Exception = archive(*source, logger)
		if except != nil {
			if *verbose {
				logger.Log(except.Error(), gopolutils.Critical)
			}
			os.Exit(1)
		}
		os.Exit(0)
	}
	var path *fayl.Path = fayl.PathFrom(*source)
	logger.Log(fmt.Sprintf("Extracting from '%s'.", path.ToString()), gopolutils.Info)
	var stripped *fayl.Path = gopolutils.Must(splitFileName(path))
	var except *gopolutils.Exception = fayl.Extract(path, stripped)
	if except != nil {
		if *verbose {
			logger.Log(except.Error(), gopolutils.Critical)
		}
	}
}
