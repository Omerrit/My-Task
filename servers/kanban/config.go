package kanban

import "gerrit-share.lan/go/utils/maps"

var typeMessages maps.String
var propMessages maps.String

func init() {
	typeMessages.Add("objobj.commentname", "[\"loopback\", 0, null, \"\"]")
	typeMessages.Add("objobj.taskname", "[\"loopback\", 0, null, \"\"]")
	typeMessages.Add("objobj.eventname", "[\"loopback\", 0, null, \"\"]")

	propMessages.Add("propobj.comment.props.commentname", "[\"loopback\", 0, \"comment\", \"\"]")
	propMessages.Add("propobj.comment.props.commentsuper", "[\"loopback\", 0, \"prop.string\", \"\"]")

	propMessages.Add("propobj.task.props.titlename", "[\"loopback\", 0, \"title\", \"\"]")
	propMessages.Add("propobj.task.props.titlesuper", "[\"loopback\", 0, \"prop.string\", \"\"]")
	propMessages.Add("propobj.task.props.end_datename", "[\"loopback\", 0, \"end_date\", \"\"]")
	propMessages.Add("propobj.task.props.end_datesuper", "[\"loopback\", 0, \"prop.date\", \"\"]")
	propMessages.Add("propobj.task.props.statusname", "[\"loopback\", 0, \"status\", \"\"]")
	propMessages.Add("propobj.task.props.statussuper", "[\"loopback\", 0, \"prop.string\", \"\"]")

	propMessages.Add("propobj.event.props.titlename", "[\"loopback\", 0, \"title\", \"\"]")
	propMessages.Add("propobj.event.props.titlesuper", "[\"loopback\", 0, \"prop.string\", \"\"]")
	propMessages.Add("propobj.event.props.start_datename", "[\"loopback\", 0, \"start_date\", \"\"]")
	propMessages.Add("propobj.event.props.start_datesuper", "[\"loopback\", 0, \"prop.date\", \"\"]")
	propMessages.Add("propobj.event.props.end_datename", "[\"loopback\", 0, \"end_date\", \"\"]")
	propMessages.Add("propobj.event.props.end_datesuper", "[\"loopback\", 0, \"prop.date\", \"\"]")
	propMessages.Add("propobj.event.props.statusname", "[\"loopback\", 0, \"status\", \"\"]")
	propMessages.Add("propobj.event.props.statussuper", "[\"loopback\", 0, \"prop.string\", \"\"]")
}
