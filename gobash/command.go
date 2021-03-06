package gobash
/* command.h -- The structures used internally to represent commands, and
   the extern declarations of the functions used to create them. */

/* Copyright (C) 1993-2009 Free Software Foundation, Inc.

   This file is part of GNU Bash, the Bourne Again SHell.

   Bash is free software: you can redistribute it and/or modify
   it under the terms of the GNU General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   Bash is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU General Public License for more details.

   You should have received a copy of the GNU General Public License
   along with Bash.  If not, see <http://www.gnu.org/licenses/>.
*/

/* Instructions describing what kind of thing to do for a redirection. */

import (
	"bytes"
	"fmt"
)

type r_instruction int

const (
  r_output_direction = r_instruction(0)
  r_input_direction = r_instruction(iota)
  r_inputa_direction = r_instruction(iota)
  r_appending_to = r_instruction(iota)
  r_reading_until = r_instruction(iota)
  r_reading_string = r_instruction(iota)
  r_duplicating_input = r_instruction(iota)
  r_duplicating_output = r_instruction(iota)
  r_deblank_reading_until = r_instruction(iota)
  r_close_this = r_instruction(iota)
  r_err_and_out = r_instruction(iota)
  r_input_output = r_instruction(iota)
  r_output_force = r_instruction(iota)
  r_duplicating_input_word = r_instruction(iota)
  r_duplicating_output_word = r_instruction(iota)
  r_move_input = r_instruction(iota)
  r_move_output = r_instruction(iota)
  r_move_input_word = r_instruction(iota)
  r_move_output_word = r_instruction(iota)
  r_append_err_and_out = r_instruction(iota)
)

/* Redirection flags; values for rflags */
const REDIR_VARASSIGN = 0x01

/* Redirection errors. */
const AMBIGUOUS_REDIRECT = -1
const NOCLOBBER_REDIRECT = -2
const RESTRICTED_REDIRECT = -3 /* can only happen in restricted shells. */
const HEREDOC_REDIRECT = -4 /* here-doc temp file can't be created */
const BADVAR_REDIRECT = -5 /* something wrong with {varname}redir */

func CLOBBERING_REDIRECT(ri r_instruction) bool {
  return ri == r_output_direction || ri == r_err_and_out
}

func OUTPUT_REDIRECT(ri r_instruction) bool {
  return (ri == r_output_direction || ri == r_input_output || ri == r_err_and_out || ri == r_append_err_and_out)
}

func INPUT_REDIRECT(ri r_instruction) bool {
  return (ri == r_input_direction || ri == r_inputa_direction || ri == r_input_output)
}

func WRITE_REDIRECT(ri r_instruction) bool {
  return (ri == r_output_direction ||
	ri == r_input_output ||
	ri == r_err_and_out ||
	ri == r_appending_to ||
	ri == r_append_err_and_out ||
	ri == r_output_force)
}

/* redirection needs translation */
func TRANSLATE_REDIRECT(ri r_instruction) bool {
  return (ri == r_duplicating_input_word || ri == r_duplicating_output_word ||
   ri == r_move_input_word || ri == r_move_output_word)
}

/* Command Types: */
type command_type int

const (
	cm_for = command_type(0)
	cm_case = command_type(iota)
	cm_while = command_type(iota)
	cm_if = command_type(iota)
	cm_simple = command_type(iota)
	cm_select = command_type(iota)
	cm_connection = command_type(iota)
	cm_function_def = command_type(iota)
	cm_until = command_type(iota)
	cm_group = command_type(iota)
	cm_arith = command_type(iota)
	cm_cond = command_type(iota)
	cm_arith_for = command_type(iota)
	cm_subshell = command_type(iota)
	cm_coproc = command_type(iota)
)

var commandTypeToName = map[command_type]string {
	cm_for: "FOR",
	cm_case: "CASE",
	cm_while: "WHILE",
	cm_if: "IF",
	cm_simple: "SIMPLE",
	cm_select: "SELECT",
	cm_connection: "CONNECTION",
	cm_function_def: "FUNCTION_DEF",
	cm_until: "UNTIL",
	cm_group: "GROUP",
	cm_arith: "ARITH",
	cm_cond: "COND",
	cm_arith_for: "ARITH_FOR",
	cm_subshell: "SUBSHELL",
	cm_coproc: "COPROC",
}

func getCommandName(typ command_type) string {
	v, ok := commandTypeToName[typ]
	if !ok {
		return fmt.Sprintf("Command typ %d is out of the enum command_type", typ)
	}
	return v
}

/* Possible values for the `flags' field of a word_desc. */
const W_HASDOLLAR = 0x000001 /* Dollar sign present. */
const W_QUOTED = 0x000002 /* Some form of quote character is present. */
const W_ASSIGNMENT = 0x000004 /* This word is a variable assignment. */
const W_GLOBEXP = 0x000008 /* This word is the result of a glob expansion. */
const W_NOSPLIT = 0x000010 /* Do not perform word splitting on this word because ifs is empty string. */
const W_NOGLOB = 0x000020 /* Do not perform globbing on this word. */
const W_NOSPLIT2 = 0x000040 /* Don't split word except for $@ expansion (using spaces) because context does not allow it. */
const W_TILDEEXP = 0x000080 /* Tilde expand this assignment word */
const W_DOLLARAT = 0x000100 /* $@ and its special handling */
const W_DOLLARSTAR = 0x000200 /* $* and its special handling */
const W_NOCOMSUB = 0x000400 /* Don't perform command substitution on this word */
const W_ASSIGNRHS = 0x000800 /* Word is rhs of an assignment statement */
const W_NOTILDE = 0x001000 /* Don't perform tilde expansion on this word */
const W_ITILDE = 0x002000 /* Internal flag for word expansion */
const W_NOEXPAND = 0x004000 /* Don't expand at all -- do quote removal */
const W_COMPASSIGN = 0x008000 /* Compound assignment */
const W_ASSNBLTIN = 0x010000 /* word is a builtin command that takes assignments */
const W_ASSIGNARG = 0x020000 /* word is assignment argument to command */
const W_HASQUOTEDNULL = 0x040000 /* word contains a quoted null character */
const W_DQUOTE = 0x080000 /* word should be treated as if double-quoted */
const W_NOPROCSUB = 0x100000 /* don't perform process substitution */
const W_HASCTLESC = 0x200000 /* word contains literal CTLESC characters */
const W_ASSIGNASSOC = 0x400000 /* word looks like associative array assignment */

/* Possible values for subshell_environment */
const SUBSHELL_ASYNC = 0x01 /* subshell caused by `command &' */
const SUBSHELL_PAREN = 0x02 /* subshell caused by ( ... ) */
const SUBSHELL_COMSUB = 0x04 /* subshell caused by `command` or $(command) */
const SUBSHELL_FORK = 0x08 /* subshell caused by executing a disk command */
const SUBSHELL_PIPE = 0x10 /* subshell from a pipeline element */
const SUBSHELL_PROCSUB = 0x20 /* subshell caused by <(command) or >(command) */
const SUBSHELL_COPROC = 0x40 /* subshell from a coproc pipeline */

/* A structure which represents a word. */
type word_desc struct {
  word string /* Zero terminated string. */
  flags int /* Flags associated with this word. */
}

func (wd *word_desc) String() string {
	if wd == nil {
		return "nil"
	}
	return fmt.Sprintf("{flags:%d, {%s}}", wd.flags, wd.word)
}

/* A linked list of words. */
type word_list struct {
  next *word_list
  word *word_desc
}

func (list *word_list) String() string {
    if list == nil {
      return "nil"
    }
	cur := list
    b := new(bytes.Buffer)
	for cur != nil {
		fmt.Fprintf(b, "%v ", cur.word)
		cur = cur.next
	}
	return string(b.Bytes())
}

func (list *word_list) Join(sep string) *StringBuilder {
  sb := NewStringBuilder()
  cur := list
  for cur != nil {
    sb.AppendString(cur.word.word)
    if cur.next != nil {
      sb.AppendString(sep)
    }
    cur = cur.next
  }
  return sb
}

func reverseWordList(list *word_list) *word_list {
	var prev *word_list

	for list != nil {
		next := list.next;
		list.next = prev;
		prev = list;
		list = next;
	}
	return prev;
}


func reverseRedirectList(list *Redirect) *Redirect {
	var prev *Redirect

	for list != nil {
		next := list.next;
		list.next = prev;
		prev = list;
		list = next;
	}
	return prev;
}

func reversePatternListList(list *PatternList) *PatternList {
	var prev *PatternList

	for list != nil {
		next := list.next;
		list.next = prev;
		prev = list;
		list = next;
	}
	return prev;
}

func makeWordList(x *word_desc, l *word_list) *word_list {
	w := new(word_list)
	w.word = x
	w.next = l
	return w
}
/* **************************************************************** */
/*								    */
/*			Shell Command Structs			    */
/*								    */
/* **************************************************************** */

/* What a redirection descriptor looks like.  If the redirection instruction
   is ri_duplicating_input or ri_duplicating_output, use DEST, otherwise
   use the file in FILENAME.  Out-of-range descriptors are identified by a
   negative DEST. */

type Redirectee struct {
  dest int /* Place to redirect REDIRECTOR to, or ... */
  filename *word_desc /* filename to redirect to. */
}

func (redir Redirectee) String() string {
  return fmt.Sprintf("dest:%d", redir.dest)
}

/* Structure describing a redirection.  If REDIRECTOR is negative, the parser
   (or translator in redir.c) encountered an out-of-range file descriptor. */
type Redirect struct {
  next *Redirect /* Next element, or NULL. */
  redirector Redirectee /* Descriptor or varname to be redirected. */
  rflags int /* Private flags for this redirection */
  flags int /* Flag value for `open'. */
  instruction r_instruction /* What to do with the information. */
  redirectee Redirectee /* File descriptor or filename */
  here_doc_eof string /* The word that appeared in <<foo. */
}

func (redir *Redirect) String() string {
  if redir == nil {
    return "nil"
  }
  cur := redir
  b := new(bytes.Buffer)
  for cur != nil {
    fmt.Fprintf(b, "{redirector:{%v}\nrflags:%d\nflags:%d\ninstruction:%v\nredirectee:{%v}}\n",
                     cur.redirector, cur.rflags, cur.flags, cur.instruction, cur.redirectee)
	cur = cur.next
  }
  return b.String()
}

/* An element used in parsing.  A single word or a single redirection.
   This is an ephemeral construct. */
type ELEMENT struct {
  word *word_desc
  redirect *Redirect
}

/* Possible values for command->flags. */
const CMD_WANT_SUBSHELL = 0x01 /* User wants a subshell: ( command ) */
const CMD_FORCE_SUBSHELL = 0x02 /* Shell needs to force a subshell. */
const CMD_INVERT_RETURN = 0x04 /* Invert the exit value. */
const CMD_IGNORE_RETURN = 0x08 /* Ignore the exit value.  For set -e. */
const CMD_NO_FUNCTIONS = 0x10 /* Ignore functions during command lookup. */
const CMD_INHIBIT_EXPANSION = 0x20 /* Do not expand the command words. */
const CMD_TIME_PIPELINE = 0x80 /* Time a pipeline */
const CMD_TIME_POSIX = 0x100 /* time -p; use POSIX.2 time output spec. */
const CMD_STDIN_REDIR = 0x400 /* async command needs implicit </dev/null */

type CommandValue struct {
    For *ForCom
    Case *CaseCom
    While *WhileCom
    If *IfCom
    Connection *Connection
    Simple *SimpleCom
    Function_def *FunctionDef
    Group *GroupCom
    Select *SelectCom
    Arith *ArithCom
    Cond *CondCom
    ArithFor *ArithForCom
    Subshell *SubshellCom
    Coproc *CoprocCom
}

type Redirectable interface {
	AddRedirect(r *Redirect)
}

/* What a command looks like. */
type Command struct {
  typ command_type /* FOR CASE WHILE IF Connection or SIMPLE. */
  flags int /* Flags controlling execution environment. */
  line int /* line number the command starts on */
  redirects *Redirect /* Special redirects for FOR CASE, etc. */
  value CommandValue
}

func (cmd *Command) String() string {
  if cmd == nil {
	return "nil"
  }
  v := cmd.value
  var out interface{}
  switch cmd.typ {
	case cm_for: out = v.For
	case cm_case: out = v.Case
	case cm_while: out = v.While
	case cm_if: out = v.If
	case cm_simple: out = v.Simple
	case cm_select: out = v.Select
	case cm_connection: out = v.Connection
	case cm_function_def: out = v.Function_def
	case cm_until: out = v.While
	case cm_group: out = v.Group
	case cm_arith: out = v.Arith
	case cm_cond: out = v.Cond
	case cm_arith_for: out = v.ArithFor
	case cm_subshell: out = v.Subshell
	case cm_coproc: out = v.Coproc
	default: panic(fmt.Sprintf("Unknown cm_type: %d", cmd.typ))
  }

  return fmt.Sprintf("{ %s{%v}}", getCommandName(cmd.typ), out)
}

func (cmd *Command) AddRedirect(r *Redirect) {
  if r == nil {
    panic("AddRedirect: r == nil")
  }
  if cmd.redirects != nil {
      var t *Redirect
      for t = cmd.redirects; t.next != nil; t = t.next {}
      t.next = r;
  } else {
    cmd.redirects = r;
  }
}

/* Structure used to represent the Connection type. */
type Connection struct {
  ignore int /* Unused; simplifies make_command (). */
  first *Command /* Pointer to the first command. */
  second *Command /* Pointer to the second command. */
  connector int /* What separates this command from others. */
}

func (conn *Connection) String() string {
  if conn == nil {
    return "nil"
  }
  return fmt.Sprintf("connector:%d\nfirst: %v\nsecond: %v",
                    conn.connector, conn.first, conn.second)
}

/* Structures used to represent the CASE command. */

/* Values for FLAGS word in a PatternList */
const CASEPAT_FALLTHROUGH = 0x01
const CASEPAT_TESTNEXT = 0x02

/* Pattern/action structure for CaseCom. */
type PatternList struct {
  next *PatternList /* Clause to try in case this one failed. */
  patterns *word_list /* Linked list of patterns to test. */
  action *Command /* Thing to execute if a pattern matches. */
  flags int
}

func (list *PatternList) String() string {
  if list == nil {
    return "nil"
  }
  b := new(bytes.Buffer)
  cur := list
  for cur != nil {
    fmt.Fprintf(b, "{patterns:%v\naction:%v\nflags:%d}\n", cur.patterns, cur.action, cur.flags)
    cur = cur.next
  }
  return b.String()
}

/* The CASE command. */
type CaseCom struct {
  flags int /* See description of CMD flags. */
  line int /* line number the `case' keyword appears on */
  word *word_desc /* The thing to test. */
  clauses *PatternList /* The clauses to test against, or NULL. */
}

func (cmd *CaseCom) String() string {
  if cmd == nil {
    return "nil"
  }
  return fmt.Sprintf("flags:%d line:%d word:%v clauses:{%v}", cmd.flags, cmd.line, cmd.word, cmd.clauses)
}

/* FOR command. */
type ForCom struct {
  flags int /* See description of CMD flags. */
  line int /* line number the `for' keyword appears on */
  name *word_desc /* The variable name to get mapped over. */
  map_list *word_list /* The things to map over.  This is never NULL. */
  action *Command	/* The action to execute.
			   During execution, NAME is bound to successive
			   members of MAP_LIST. */
}

type ArithForCom struct {
  flags int
  line int /* generally used for error messages */
  init *word_list
  test *word_list
  step *word_list
  action *Command
}

/* KSH SELECT command. */
type SelectCom struct {
  flags int /* See description of CMD flags. */
  line int /* line number the `select' keyword appears on */
  name *word_desc /* The variable name to get mapped over. */
  map_list *word_list /* The things to map over.  This is never NULL. */
  action *Command	/* The action to execute.
			   During execution, NAME is bound to the member of
			   MAP_LIST chosen by the user. */
}

/* IF command. */
type IfCom struct {
  flags int /* See description of CMD flags. */
  test *Command /* Thing to test. */
  true_case *Command /* What to do if the test returned non-zero. */
  false_case *Command /* What to do if the test returned zero. */
}

func (cmd *IfCom) String() string {
  if cmd == nil {
    return "nil"
  }
  return fmt.Sprintf("flags:%d\ntest:%v\ntrue_case:%v\nfalse_case:%v", cmd.flags, cmd.test, cmd.true_case, cmd.false_case);
}

/* WHILE command. */
type WhileCom struct {
  flags int /* See description of CMD flags. */
  test *Command /* Thing to test. */
  action *Command /* Thing to do while test is non-zero. */
}

/* The arithmetic evaluation command, ((...)).  Just a set of flags and
   a word_list, of which the first element is the only one used, for the
   time being. */
type ArithCom struct {
  flags int
  line int
  exp *word_list
}

/* The conditional command, [[...]].  This is a binary tree -- we slippped
   a recursive-descent parser into the YACC grammar to parse it. */
const COND_AND = 1
const COND_OR = 2
const COND_UNARY = 3
const COND_BINARY = 4
const COND_TERM = 5
const COND_EXPR = 6

type CondCom struct {
  flags int
  line int
  typ int
  op *word_desc
  left *CondCom
  right *CondCom
}

func (cmd *CondCom) String() string {
  if cmd == nil {
    return "nil"
  }
  return fmt.Sprintf("%d: flags:%d\ntyp:%d\nop:%v\nleft:{%v}\nright:{%v}", cmd.line, cmd.flags, cmd.typ, cmd.op, cmd.left, cmd.right)
}

/* The "simple" command.  Just a collection of words and redirects. */
type SimpleCom struct {
  flags int /* See description of CMD flags. */
  line int /* line number the command starts on */
  words *word_list;		/* The program name, the arguments,
				   variable assignments, etc. */
  redirects *Redirect /* Redirections to perform. */
}

func (cmd *SimpleCom) String() string {
	if cmd == nil {
		return "nil"
	}
    return fmt.Sprintf( "%d: words:%vredirects:{%v}", cmd.line, cmd.words, cmd.redirects)
}

func (cmd *SimpleCom) AddRedirect(r *Redirect) {
  if r == nil {
    panic("AddRedirect: r == nil")
  }
  if cmd.redirects != nil {
      var t *Redirect
      for t = cmd.redirects; t.next != nil; t = t.next {}
      t.next = r;
  } else {
    cmd.redirects = r;
  }
}


/* The "function definition" command. */
type FunctionDef struct {
  flags int /* See description of CMD flags. */
  line int /* Line number the function def starts on. */
  name *word_desc /* The name of the function. */
  command *Command /* The parsed execution tree. */
  source_file string /* file in which function was defined, if any */
}

/* A command that is `grouped' allows pipes and redirections to affect all
   commands in the group. */
type GroupCom struct {
  ignore int /* See description of CMD flags. */
  command *Command
}

type SubshellCom struct {
  flags int
  command *Command
}

func (cmd *SubshellCom) String() string {
  if cmd == nil {
    return "nil"
  }
  return fmt.Sprintf("flags:%d\ncommand: %v", cmd.flags, cmd.command)
}

const COPROC_RUNNING = 0x01
const COPROC_DEAD = 0x02

type coproc struct {
  c_name string
  c_pid uint32
  c_rfd int
  c_wfd int
  c_rsave int
  c_wsave int
  c_flags int
  c_status int
}

type CoprocCom struct {
  flags int
  name string
  command *Command
}

/* Possible command errors */
const CMDERR_DEFAULT = 0
const CMDERR_BADTYPE = 1
const CMDERR_BADCONN = 2
const CMDERR_BADJUMP = 3

const CMDERR_LAST = 3

